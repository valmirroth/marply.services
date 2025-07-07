import { useState, useEffect, useCallback } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { apiRequest } from "@/lib/queryClient";
import { insertInventoryCountSchema, type InsertInventoryCount } from "@shared/schema";
import { z } from "zod";

const formSchema = insertInventoryCountSchema.extend({
  itemCode: z.string().min(1, "Código do item é obrigatório"),
  location: z.string().min(1, "Local é obrigatório"),
  quantity: z.string().min(1, "Quantidade é obrigatória"),
  volumeCount: z.string().min(0.001, "Volume é obrigatório"),
  empresa: z.number().min(1, "Empresa é obrigatória"),
});

type FormData = z.infer<typeof formSchema>;

interface InventoryFormProps {
  onSuccess: () => void;
  disabled?: boolean;
}

export default function InventoryForm({ onSuccess, disabled = false }: InventoryFormProps) {
  const [itemCode, setItemCode] = useState("");
  const [itemDescription, setItemDescription] = useState("");
  const [isLoadingDescription, setIsLoadingDescription] = useState(false);
  const [isItemFound, setIsItemFound] = useState(true);
  const [location, setLocation] = useState("");
  const [locationInfo, setLocationInfo] = useState<{ codigo: string, sigla: string, descricao: string } | null>(null);
  const [isLocationValid, setIsLocationValid] = useState(true);
  const [isLoadingLocation, setIsLoadingLocation] = useState(false);
  const [isApiAvailable, setIsApiAvailable] = useState(true);
  const [qtdeEmbalagem, setQtdeEmbalagem] = useState<number>(0);
  const [hasQuantityAlert, setHasQuantityAlert] = useState(false);
  const [isCalculateAuto, setCalculateAuto] = useState(false);
  const [empresa, setEmpresa] = useState<number>(() => {
    const savedEmpresa = localStorage.getItem("selectedEmpresa");
    return savedEmpresa ? parseInt(savedEmpresa) : 1;
  });
  const { toast } = useToast();

  // Executar validação quando qtdeEmbalagem mudar
  useEffect(() => {
    if (qtdeEmbalagem > 0) {
      validarEmbalagem();
    }
  }, [qtdeEmbalagem]);

  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      itemCode: "",
      itemDescription: "",
      location: "",
      quantity: "",
      volumeCount: "",
      empresa: empresa,
    },
  });
  // Função para alterar empresa
  const handleEmpresaChange = (novaEmpresa: number) => {
    setEmpresa(novaEmpresa);
    localStorage.setItem("selectedEmpresa", novaEmpresa.toString());
    form.setValue("empresa", novaEmpresa);
  };

  // Função para buscar descrição do item
  const buscarDescricaoItem = async (codigo: string) => {
    if (!codigo.trim()) {
      setItemDescription("");
      setQtdeEmbalagem(0);
      setIsItemFound(true);
      return;
    }

    setIsLoadingDescription(true);
    setIsItemFound(true);

    // Dados simulados com qtde_embalagem para demonstração
    const itensSimulados = [
      { codigo: "1", descricao: "TESTE DO CÓDIGO", qtde_embalagem: 10 },
      { codigo: "A001234", descricao: "Parafuso Sextavado M12x50 Aço Inox", qtde_embalagem: 25 },
      { codigo: "B002456", descricao: "Arruela Lisa M12 Aço Galvanizado", qtde_embalagem: 50 },
      { codigo: "C003789", descricao: "Porca Sextavada M12 Aço Inox", qtde_embalagem: 20 },
      { codigo: "D004567", descricao: "Parafuso Phillips M8x40 Aço Carbono", qtde_embalagem: 15 },
      { codigo: "E005123", descricao: "Porca Flangeada M10 Aço Inox", qtde_embalagem: 30 }
    ];

    // Simular delay de API
    await new Promise(resolve => setTimeout(resolve, 800));

    let itemEncontrado = false;

    try {
      // Tentar buscar na API Go primeiro
      const response = await fetch(`http://localhost:9090/api/items/${codigo}/description`);
      const data = await response.json();

      if (data.success && data.data && data.data.descricao !== "Descrição não encontrada") {
        setItemDescription(data.data.descricao);
        form.setValue("itemDescription", data.data.descricao);
        setQtdeEmbalagem(data.data.qtde_embalagem || 0);
        setCalculateAuto(data.data.calcular_automatico || false);
        setIsItemFound(true);
        itemEncontrado = true;
        setIsLoadingDescription(false);
        return;
      }
    } catch (error) {
      // Se API não responder, usar dados simulados
    }

    // Usar dados simulados como fallback
    const itemSimulado = itensSimulados.find(item =>
      item.codigo.toLowerCase() === codigo.toLowerCase()
    );

    if (itemSimulado) {
      setItemDescription(itemSimulado.descricao);
      form.setValue("itemDescription", itemSimulado.descricao);
      setQtdeEmbalagem(itemSimulado.qtde_embalagem);
      setIsItemFound(true);
      itemEncontrado = true;
    } else {
      setItemDescription("Descrição não encontrada");
      form.setValue("itemDescription", "Descrição não encontrada");
      setQtdeEmbalagem(0);
      setIsItemFound(false);

      // Emitir toast quando item não for encontrado
      toast({
        title: "Item não encontrado",
        description: `O código "${codigo}" não foi encontrado na base de dados`,
        variant: "destructive",
      });
    }

    setIsLoadingDescription(false);
  };

  // Função para validar local
  const validarLocal = async (sigla: string) => {
    if (!sigla.trim()) {
      setLocationInfo(null);
      setIsLocationValid(true);
      setIsApiAvailable(true);
      return;
    }

    setIsLoadingLocation(true);
    try {
      const emp = localStorage.getItem("selectedEmpresa");
      const response = await fetch(`http://localhost:9090/api/locations/${sigla}/${emp}/validate`);

      setIsApiAvailable(true);
      const data = await response.json();

      if (data.success && data.data) {
        setLocationInfo(data.data);
        setIsLocationValid(true);
        toast({
          title: "Local válido",
          description: `${data.data.descricao} (${data.data.codigo})`,
        });
      } else {
        setLocationInfo(null);
        setIsLocationValid(false);
        toast({
          title: "Local inválido",
          description: "Local não encontrado na base de dados",
          variant: "destructive",
        });
      }
    } catch (error) {
      setLocationInfo(null);
      setIsLocationValid(false);
      setIsApiAvailable(false);

      toast({
        title: "API não disponível",
        description: "Não foi possível comunicar com a API. Verifique a conexão.",
        variant: "destructive",
      });
    } finally {
      setIsLoadingLocation(false);
    }
  };

  // Função para validar embalagem
  const validarEmbalagem = useCallback(() => {

    const volumes = parseFloat(form.getValues("volumeCount")) || 0.000;
    if (isCalculateAuto) {
      const volumes = parseFloat(form.getValues("volumeCount")) || 0.000;
      const x = (volumes * qtdeEmbalagem);
      form.setValue("quantity", x.toFixed(3));
      console.log("setor o campo valor qtde.")
    }
    const quantidade = parseFloat(form.getValues("quantity")) || 0.000;

    if (qtdeEmbalagem > 0 && quantidade > 0 && volumes > 0) {
      const quantidadeCalculada = volumes * qtdeEmbalagem;
      const diferenca = Math.abs(quantidade - quantidadeCalculada);

      if (diferenca > 1) {
        setHasQuantityAlert(true);
        toast({
          title: "Conferir quantidades",
          //  description: `Quantidade esperada: ${quantidadeCalculada.toFixed(5)} (${volumes} × ${qtdeEmbalagem}). Diferença: ${diferenca.toFixed(5)}`,
          variant: "destructive",
        });
      } else {
        setHasQuantityAlert(false);
      }
    } else {
      setHasQuantityAlert(false);
    }
  }, [qtdeEmbalagem, form, toast]);

  const createCountMutation = useMutation({
    mutationFn: async (data: InsertInventoryCount) => {
      const response = await apiRequest("POST", "/api/contagens", data);
      return response.json();
    },
    onSuccess: () => {
      toast({
        title: "Sucesso",
        description: "Contagem registrada com sucesso!",
      });
      form.reset({
        itemCode: "",
        itemDescription: "",
        location: "",
        quantity: "",
        volumeCount: "",
        empresa: empresa, // Manter empresa selecionada
      });
      form.reset();
      setItemCode("");
      setItemDescription("");
      setIsItemFound(true);
      setLocation("");
      setLocationInfo(null);
      setIsLocationValid(true);
      setIsApiAvailable(true);
      setQtdeEmbalagem(0);
      setHasQuantityAlert(false);
      onSuccess();
    },
    onError: (error: any) => {
      const errorMessage = error?.message?.includes("fetch")
        ? "API não disponível. Verifique a conexão."
        : error?.message || "Erro ao registrar contagem";

      toast({
        title: "Erro",
        description: errorMessage,
        variant: "destructive",
      });
    },
  });

  const onSubmit = (data: FormData) => {
    const submitData: InsertInventoryCount = {
      itemCode: data.itemCode,
      itemDescription: itemDescription || "Descrição não encontrada",
      location: data.location,
      quantity: data.quantity,
      volumeCount: data.volumeCount,
      empresa: data.empresa,
    };

    createCountMutation.mutate(submitData);
  };

  const handleItemCodeBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const code = e.target.value.trim();
    setItemCode(code);
    buscarDescricaoItem(code);
  };

  const handleLocationBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const sigla = e.target.value.trim().toUpperCase();
    setLocation(sigla);
    validarLocal(sigla);
  };

  if (disabled) {
    return null;
  }

  return (
    <Card className="shadow-lg">
      <CardContent className="p-8">
        <div className="mb-6">
          <h2 className="text-xl font-semibold text-gray-900 flex items-center">
            <svg className="w-5 h-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Nova Contagem
          </h2>
        </div>

        {/* Campo de Empresa */}
        <div className="mt-4 p-4 bg-gray-50 rounded-lg">
          <Label className="text-sm font-medium flex items-center mb-3">
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h4a1 1 0 011 1v5m-6 0V9a1 1 0 011-1h4a1 1 0 011 1v11" />
            </svg>
            Empresa *
          </Label>
          <div className="flex gap-4">
            <button
              type="button"
              onClick={() => handleEmpresaChange(1)}
              className={`flex-1 p-3 rounded-lg border-2 transition-all duration-200 ${empresa === 1
                ? "border-blue-500 bg-blue-50 text-blue-700"
                : "border-gray-300 bg-white text-gray-700 hover:border-gray-400"
                }`}
            >
              <div className="flex items-center justify-center">
                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h4a1 1 0 011 1v5m-6 0V9a1 1 0 011-1h4a1 1 0 011 1v11" />
                </svg>
                Matriz
              </div>
            </button>
            <button
              type="button"
              onClick={() => handleEmpresaChange(5)}
              className={`flex-1 p-3 rounded-lg border-2 transition-all duration-200 ${empresa === 5
                ? "border-blue-500 bg-blue-50 text-blue-700"
                : "border-gray-300 bg-white text-gray-700 hover:border-gray-400"
                }`}
            >
              <div className="flex items-center justify-center">
                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 14v3m4-3v3m4-3v3M3 21h18M3 10h18M3 7l9-4 9 4M4 10h16v11H4V10z" />
                </svg>
                Filial
              </div>
            </button>
          </div>
        </div>

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="itemCode" className="text-sm font-medium flex items-center">
              <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />
              </svg>
              Código do Item *
            </Label>
            <Input
              id="itemCode"
              data-testid="input-itemCode"
              {...form.register("itemCode")}
              onBlur={handleItemCodeBlur}
              className="h-12 text-base"
              placeholder="Digite o código do item"
            />
            {form.formState.errors.itemCode && (
              <p className="text-red-500 text-sm">{form.formState.errors.itemCode.message}</p>
            )}
          </div>

          {/* Campo de Descrição do Item */}
          <div className="space-y-2">
            <Label htmlFor="itemDescription" className="text-sm font-medium flex items-center">
              <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Descrição do Item
            </Label>
            <div className="relative">
              <Input
                id="itemDescription"
                data-testid="input-item-description"
                value={itemDescription}
                className="h-12 text-base bg-gray-50"
                placeholder="A descrição será preenchida automaticamente"
                disabled={true}
              />
              {isLoadingDescription && (
                <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
                </div>
              )}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="location" className="text-sm font-medium flex items-center">
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 616 0z" />
                </svg>
                Local da Contagem *
              </Label>
              <div className="relative">
                <Input
                  id="location"
                  data-testid="input-location"
                  {...form.register("location")}
                  className={`h-12 text-base ${!isLocationValid ? "border-red-500 bg-red-50" : ""
                    }`}
                  placeholder="Ex: EST-A01"
                  onBlur={handleLocationBlur}
                  disabled={createCountMutation.isPending}
                />
                {isLoadingLocation && (
                  <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
                  </div>
                )}
              </div>
              {!isLocationValid && location && (
                <p className="text-red-500 text-sm flex items-center mt-1">
                  <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                  </svg>
                  {!isApiAvailable ? "API não disponível" : "Local inválido"} - os campos abaixo estão desabilitados
                </p>
              )}
              {locationInfo && (
                <p className="text-green-600 text-sm flex items-center mt-1">
                  <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                  </svg>
                  {locationInfo.descricao} ({locationInfo.codigo})
                </p>
              )}
              {form.formState.errors.location && (
                <p className="text-red-500 text-sm">{form.formState.errors.location.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="quantity" className="text-sm font-medium flex items-center">
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
                </svg>
                Quantidade (5 casas decimais) *
              </Label>
              <Input
                id="quantity"
                data-testid="input-quantity"
                type="number"
                step="0.00001"
                min="0"
                {...form.register("quantity", {
                  onChange: validarEmbalagem
                })}
                className={`h-12 text-base ${!isLocationValid && location ? "bg-gray-100 cursor-not-allowed" : ""
                  } ${hasQuantityAlert ? "border-red-500 bg-red-50" : ""}`}
                placeholder="0.00000"
                disabled={(!isLocationValid && location !== "") || createCountMutation.isPending || isCalculateAuto}
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
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="volumeCount" className="text-sm font-medium flex items-center">
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
                Quantidade de Volumes *
              </Label>
              <Input
                id="volumeCount"
                data-testid="input-volumeCount"
                type="number"
                min="0.000"
                step="0.001"
                {...form.register("volumeCount", {
                  onChange: validarEmbalagem
                })}
                className={`h-12 text-base ${!isLocationValid && location ? "bg-gray-100 cursor-not-allowed" : ""
                  } ${hasQuantityAlert ? "border-red-500 bg-red-50" : ""}`}
                placeholder="0.000"
                disabled={(!isLocationValid && location !== "") || createCountMutation.isPending}
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

          <Button
            type="submit"
            data-testid="button-submit"
            disabled={(!isLocationValid && location !== "") || createCountMutation.isPending || !isItemFound || hasQuantityAlert}
            className="w-full h-12 text-base font-semibold bg-blue-600 hover:bg-blue-700"
          >
            {createCountMutation.isPending ? (
              <div className="flex items-center">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                Registrando...
              </div>
            ) : (
              "Registrar Contagem"
            )}
          </Button>
        </form>
      </CardContent>
    </Card >
  );
}