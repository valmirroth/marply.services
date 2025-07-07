import { useMutation } from "@tanstack/react-query";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { useToast } from "@/hooks/use-toast";
import { apiRequest } from "@/lib/queryClient";
import type { InventoryCount } from "@shared/schema";

interface DeleteModalProps {
  count: InventoryCount;
  onClose: () => void;
  onSuccess: () => void;
}

export default function DeleteModal({ count, onClose, onSuccess }: DeleteModalProps) {
  const { toast } = useToast();

  const deleteCountMutation = useMutation({
    mutationFn: async () => {
      const response = await apiRequest("DELETE", `http://localhost:9090/api/inventory-counts/${count.id}`);
      return response.json();
    },
    onSuccess: () => {
      toast({
        title: "Sucesso",
        description: "Contagem excluída com sucesso!",
      });
      onSuccess();
    },
    onError: (error: any) => {
      toast({
        title: "Erro",
        description: error.message || "Erro ao excluir contagem",
        variant: "destructive",
      });
    },
  });

  const handleConfirmDelete = () => {
    deleteCountMutation.mutate();
  };

  return (
    <Dialog open={true} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center text-red-600">
            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            Confirmar Exclusão
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <p>Tem certeza que deseja excluir esta contagem?</p>

          <Alert className="border-yellow-300 bg-yellow-50">
            <AlertDescription>
              <div className="space-y-1">
                <div><strong>Código:</strong> {count.itemCode}</div>
                <div><strong>Local:</strong> {count.location}</div>
                <div><strong>Quantidade:</strong> {count.quantity}</div>
              </div>
            </AlertDescription>
          </Alert>

          <p className="text-sm text-gray-600">Esta ação não pode ser desfeita.</p>
        </div>

        <div className="flex justify-end gap-3 pt-4">
          <Button type="button" variant="outline" onClick={onClose} data-testid="button-cancel-delete">
            Cancelar
          </Button>
          <Button
            variant="destructive"
            onClick={handleConfirmDelete}
            disabled={deleteCountMutation.isPending}
            data-testid="button-confirm-delete"
          >
            {deleteCountMutation.isPending ? "Excluindo..." : "Excluir"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
