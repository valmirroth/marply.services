import { useState, useEffect } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import InventoryForm from "@/components/inventory-form";
import InventoryList from "@/components/inventory-list";
import EditModal from "@/components/edit-modal";
import DeleteModal from "@/components/delete-modal";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useToast } from "@/hooks/use-toast";
import { apiRequest } from "@/lib/queryClient";
import type { InventoryCount } from "@shared/schema";

export default function InventoryPage() {
  const [editingCount, setEditingCount] = useState<InventoryCount | null>(null);
  const [deletingCount, setDeletingCount] = useState<InventoryCount | null>(null);
  const [isCountFinalized, setIsCountFinalized] = useState(false);
  const [empresaSelecionada, setEmpresaSelecionada] = useState<string>("");
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const savedEmpresa = localStorage.getItem("selectedEmpresa");
  const { data: counts = [], refetch, isLoading, error } = useQuery<InventoryCount[]>({
    queryKey: ["/api/inventory-counts", localStorage.getItem("selectedEmpresa")],
  });


  // Query para verificar status da contagem
  const { data: statusData } = useQuery({
    queryKey: ["contagem-status"],
    queryFn: async () => {
      try {
        const _empresa = localStorage.getItem("selectedEmpresa");
        // Tentar API Go primeiro
        const response = await fetch("http://localhost:9090/api/contagensstate/status/" + _empresa);
        if (response.ok) {
          return response.json();
        }
      } catch (error) {
        // Se API Go não responder, retornar estado padrão
        console.log("API Go não disponível, usando estado padrão");
      }

      // Fallback: assumir não finalizada se API não responder
      return { success: true, data: { finalizada: false } };
    },
    refetchInterval: 5000, // Verifica a cada 5 segundos
  });

  // Mutation para finalizar contagem
  const finalizarContagemMutation = useMutation({
    mutationFn: async () => {
      try {
        // Tentar API Go primeiro
        const _empresa = localStorage.getItem("selectedEmpresa");
        const response = await fetch("http://localhost:9090/api/contagens/finalizar/" + _empresa, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
        });

        if (response.ok) {
          const data = await response.json();
          if (data.success) {
            return data;
          } else {
            throw new Error(data.error || "Erro ao finalizar contagem");
          }
        } else {
          throw new Error(`Erro HTTP: ${response.status}`);
        }
      } catch (error) {
        console.error("Erro ao finalizar contagem:", error);
        throw error;
      }
    },
    onSuccess: () => {
      setIsCountFinalized(true);
      refetch();
      // Invalidar cache da query de status para forçar nova consulta
      queryClient.invalidateQueries({ queryKey: ["contagem-status"] });
      toast({
        title: "Contagem finalizada",
        description: "A contagem foi encerrada com sucesso!",
      });
    },
    onError: (error: any) => {
      toast({
        title: "Erro ao finalizar",
        description: error?.message || "Erro ao finalizar contagem",
        variant: "destructive",
      });
    },
  });

  // Verificar status ao carregar e quando statusData mudar
  useEffect(() => {
    if (statusData?.success && statusData?.data?.finalizada) {
      setIsCountFinalized(true);
    } else if (statusData?.success && !statusData?.data?.finalizada) {
      setIsCountFinalized(false);
    }
  }, [statusData]);

  // Garantir que counts seja sempre um array
  const safeCounts = Array.isArray(counts) ? counts : [];

  const currentDate = new Date().toLocaleDateString('pt-BR');

  return (
    <div className="min-h-screen bg-slate-50">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-blue-800 text-white sticky top-0 z-50 shadow-lg">
        <div className="container mx-auto px-4 py-6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-semibold mb-1 flex items-center">
                <svg className="w-6 h-6 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
                </svg>
                Inventário
              </h1>
              <p className="text-blue-100 text-sm">Sistema móvel para contagem de estoque</p>
            </div>
            <div className="text-right text-sm">
              <div className="text-blue-100">Usuário: Admin</div>
              <div className="text-blue-100">{currentDate}</div>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6 space-y-6">
        {isCountFinalized && (
          <Alert className="border-green-300 bg-green-50">
            <AlertDescription className="flex items-center text-green-800">
              <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <strong>Contagem Finalizada:</strong> A contagem foi encerrada. Não é possível adicionar, editar ou excluir mais itens.
            </AlertDescription>
          </Alert>
        )}

        <InventoryForm
          onSuccess={() => refetch()}
          disabled={isCountFinalized}

        />

        <InventoryList
          counts={safeCounts}
          onEdit={isCountFinalized ? () => { } : setEditingCount}
          onDelete={isCountFinalized ? () => { } : setDeletingCount}
          disabled={isCountFinalized}
        />

        {!isCountFinalized && safeCounts.length > 0 && (
          <div className="flex justify-center pt-4">
            <Button
              onClick={() => finalizarContagemMutation.mutate()}
              data-testid="button-finalize-count"
              disabled={finalizarContagemMutation.isPending}
              className="bg-red-600 hover:bg-red-700 text-white px-8 py-3 text-lg font-semibold"
            >
              {finalizarContagemMutation.isPending ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Finalizando...
                </div>
              ) : (
                <>
                  <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Encerrar Contagem
                </>
              )}
            </Button>
          </div>
        )}
      </div>

      {/* Floating Summary */}
      <div className="fixed bottom-5 right-5 bg-blue-600 text-white rounded-full px-6 py-3 shadow-lg z-40">
        <div className="flex items-center text-sm font-semibold">
          <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          {safeCounts.length} Contagens
        </div>
      </div>

      {editingCount && (
        <EditModal
          isOpen={!!editingCount}
          onClose={() => setEditingCount(null)}
          item={editingCount}
        />
      )}

      {deletingCount && (
        <DeleteModal
          count={deletingCount}
          onClose={() => setDeletingCount(null)}
          onSuccess={() => {
            refetch();
            setDeletingCount(null);
          }}
        />
      )}
    </div>
  );
}
