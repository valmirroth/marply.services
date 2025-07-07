import { useEffect, useCallback, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { apiRequest } from "@/lib/queryClient";
import { updateInventoryCountSchema, type UpdateInventoryCount, type InventoryCount } from "@shared/schema";

interface EditModalProps {
  isOpen: boolean;
  onClose: () => void;
  item: InventoryCount | null;
}

export default function EditModal({ isOpen, onClose, item }: EditModalProps) {
  const [isLocationValid, setIsLocationValid] = useState(true);
  const [isLoadingLocation, setIsLoadingLocation] = useState(false);
  const [locationInfo, setLocationInfo] = useState<{ codigo: string, sigla: string, descricao: string } | null>(null);
  const [qtdeEmbalagem, setQtdeEmbalagem] = useState<number>(0);
  const [hasQuantityAlert, setHasQuantityAlert] = useState(false);
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const form = useForm<UpdateInventoryCount>({
    resolver: zodResolver(updateInventoryCountSchema),
    defaultValues: {
      itemCode: "",
      itemDescription: "",
      location: "",
      quantity: "",
      volumeCount: 0,
    },
  });

  // Função para validar embalagem
  const validarEmbalagem = useCallback(() => {
    const quantidade = parseFloat(form.getValues("quantity")) || 0;
    const volumes = parseInt(form.getValues("volumeCount")?.toString()) || 0;
    console.log("Quantidade:", quantidade);
    console.log("Volumes:", volumes);
    console.log("Qtde Embalagem:", qtdeEmbalagem);
    if (qtdeEmbalagem > 0 && quantidade > 0 && volumes > 0) {
      const quantidadeCalculada = volumes * qtdeEmbalagem;
      const diferenca = Math.abs(quantidade - quantidadeCalculada);

      if (diferenca > 1) {
        setHasQuantityAlert(true);
        toast({
          title: "Conferir quantidades",
          // description: `Quantidade esperada: ${quantidadeCalculada.toFixed(0)} (${volumes} × ${qtdeEmbalagem}). Diferença: ${diferenca.toFixed(0)}`,
          variant: "destructive",
        });
      } else {
        setHasQuantityAlert(false);
      }
    } else {
      setHasQuantityAlert(false);
    }
  }, [qtdeEmbalagem, form, toast]);



  //FIM do código

  useEffect(() => {
    if (item && isOpen) {
      form.setValue("itemCode", item.itemCode);
      form.setValue("itemDescription", item.itemDescription);
      form.setValue("location", item.location);
      form.setValue("quantity", item.quantity);
      form.setValue("volumeCount", item.volumeCount);

      // Validar o local quando carregar os dados
      validarLocal(item.location);

      // Buscar quantidade de embalagem do item
      buscarQuantidadeEmbalagem(item.itemCode);
    }
  }, [item, isOpen]);

  // Executar validação quando qtdeEmbalagem mudar
  useEffect(() => {
    if (qtdeEmbalagem > 0) {
      validarEmbalagem();
    }
  }, [qtdeEmbalagem]);

  // Função para buscar quantidade de embalagem do item
  const buscarQuantidadeEmbalagem = async (codigo: string) => {
    if (!codigo.trim()) {
      setQtdeEmbalagem(0);
      return;
    }

    try {
      const response = await fetch(`http://localhost:9090/api/items/${codigo}/description`);
      const data = await response.json();

      if (data.success && data.data) {
        setQtdeEmbalagem(data.data.qtde_embalagem || 0);
      } else {
        setQtdeEmbalagem(0);
      }
    } catch (error) {
      setQtdeEmbalagem(0);
    }
  };

  const updateMutation = useMutation({
    mutationFn: async (data: UpdateInventoryCount) => {
      const response = await apiRequest("PUT", `http://localhost:9090/api/inventory-counts/${item.id}`, data);
      return response.json();
    },
    onSuccess: () => {
      toast({
        title: "Sucesso",
        description: "Contagem atualizada com sucesso!",
      });
      queryClient.invalidateQueries({ queryKey: ["inventory-counts"] });
      onClose();
    },
    onError: (error: any) => {
      toast({
        title: "Erro",
        description: error.message || "Erro ao atualizar contagem",
        variant: "destructive",
      });
    },
  });

  const onSubmit = (data: UpdateInventoryCount) => {
    if (hasQuantityAlert) {
      toast({
        title: "Verificação pendente",
        description: "Por favor, corrija as inconsistências de quantidade e embalagem antes de salvar.",
        variant: "destructive",
      });
      return;
    }
    updateMutation.mutate(data);
  };

  // Função para validar local
  const validarLocal = async (sigla: string) => {
    if (!sigla.trim()) {
      setLocationInfo(null);
      setIsLocationValid(true);
      return;
    }

    setIsLoadingLocation(true);
    try {
      const response = await fetch(`http://localhost:9090/api/locations/${sigla}/validate`);
      const data = await response.json();

      if (data.success && data.data) {
        setLocationInfo(data.data);
        setIsLocationValid(true);
      } else {
        setLocationInfo(null);
        setIsLocationValid(false);
      }
    } catch (error) {
      setLocationInfo(null);
      setIsLocationValid(false);
    } finally {
      setIsLoadingLocation(false);
    }
  };

  const handleLocationChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newLocation = e.target.value;
    form.setValue("location", newLocation);
    validarLocal(newLocation);
  };

  const handleClose = () => {
    form.reset();
    setIsLocationValid(true);
    setIsLoadingLocation(false);
    setLocationInfo(null);
    setQtdeEmbalagem(0);
    setHasQuantityAlert(false);
    onClose();
  };

  if (!item) {
    return null;
  }

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center">
            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Editar Contagem de Item: {item.itemCode}
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="itemCode">Código do Item</Label>
            <Input
              id="itemCode"
              value={item.itemCode}
              disabled
              className="bg-gray-50"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="location">Local</Label>
            <Input
              id="location"
              data-testid="input-edit-location"

              {...form.register("location", { onChange: handleLocationChange })}
              className={isLocationValid ? "" : "border-red-500 bg-red-50"}
              disabled
            />
            {isLoadingLocation && <p className="text-sm text-gray-500">Verificando local...</p>}
            {!isLocationValid && !isLoadingLocation && (
              <p className="text-red-500 text-sm">Local inválido.</p>
            )}
            {locationInfo && (
              <div className="text-sm text-gray-600">
                <p>Código: {locationInfo.codigo}</p>
                <p>Sigla: {locationInfo.sigla}</p>
                <p>Descrição: {locationInfo.descricao}</p>
              </div>
            )}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="quantity">Quantidade (5 casas decimais)</Label>
              <Input
                id="quantity"
                {...form.register("quantity", {
                  onChange: validarEmbalagem
                })}
                type="number"
                step="0.00001"
                min="0"
                className={hasQuantityAlert ? "border-red-500 bg-red-50" : ""}
                disabled={updateMutation.isPending}
              />
              {form.formState.errors.quantity && (
                <p className="text-red-500 text-sm">{form.formState.errors.quantity.message}</p>
              )}
              {hasQuantityAlert && (
                <p className="text-red-500 text-sm flex items-center mt-1">
                  <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                  </svg>
                  Verifique a quantidade - diferença de embalagem detectada
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="volumeCount">Quantidade de Volumes</Label>
              <Input
                id="volumeCount"
                {...form.register("volumeCount", {
                  valueAsNumber: true,
                  onChange: validarEmbalagem
                })}
                type="number"
                min="0"
                step="1"
                className={hasQuantityAlert ? "border-red-500 bg-red-50" : ""}
                disabled={updateMutation.isPending}
              />
              {form.formState.errors.volumeCount && (
                <p className="text-red-500 text-sm">{form.formState.errors.volumeCount.message}</p>
              )}
              {hasQuantityAlert && (
                <p className="text-red-500 text-sm flex items-center mt-1">
                  <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                  </svg>
                  Verifique os volumes - diferença de embalagem detectada
                </p>
              )}
            </div>
          </div>

          <div className="flex justify-end gap-3 pt-4">
            <Button type="button" variant="outline" onClick={handleClose} data-testid="button-cancel-edit">
              Cancelar
            </Button>
            <Button
              type="submit"
              disabled={updateMutation.isPending || !isLocationValid || hasQuantityAlert}
              data-testid="button-save-edit"
            >
              {updateMutation.isPending ? "Salvando..." : "Salvar Alterações"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}