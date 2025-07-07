import { QueryClient, QueryFunction } from "@tanstack/react-query";

// URL base da API Go
const API_BASE_URL = "http://localhost:9090";

async function throwIfResNotOk(res: Response) {
  if (!res.ok) {
    const text = (await res.text()) || res.statusText;
    throw new Error(`${res.status}: ${text}`);
  }
}

export async function apiRequest(
  method: string,
  url: string,
  data?: unknown | undefined,
): Promise<Response> {
  // Construir URL completa para a API Go
  const fullUrl = url.startsWith("http") ? url : `${API_BASE_URL}${url}`;

  const res = await fetch(fullUrl, {
    method,
    headers: data ? { "Content-Type": "application/json" } : {},
    body: data ? JSON.stringify(data) : undefined,
    //credentials: "include",
  });

  await throwIfResNotOk(res);
  return res;
}

// Função específica para requisições que esperam resposta da API Go
export async function apiRequestGo(
  method: string,
  url: string,
  data?: unknown | undefined,
): Promise<any> {
  const res = await apiRequest(method, url, data);
  const jsonData = await res.json();

  // Se a resposta tem estrutura {success, message, data}, extrair apenas data
  if (
    jsonData &&
    typeof jsonData === "object" &&
    "success" in jsonData &&
    "data" in jsonData
  ) {
    return jsonData.data;
  }

  return jsonData;
}

type UnauthorizedBehavior = "returnNull" | "throw";
// Função para limpar e validar GUID do SQL Server
function cleanSqlServerGuid(rawId: any): string {
  if (!rawId) {
    return `temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  // Se for string, tentar limpar caracteres especiais
  if (typeof rawId === "string") {
    // Remover caracteres não-ASCII e manter apenas caracteres válidos para GUID
    const cleaned = rawId.replace(/[^\w\-]/g, "");
    if (cleaned.length >= 8) {
      return cleaned;
    }
  }

  // Se não conseguir processar, gerar um ID temporário único
  return `item-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

// Função para processar dados da API 
function processApiResponse(data: any[], url: string) {
  // Se os dados vêm da API Go (com campos como codigo_item, descricao_item)
  if (url.includes("/api/contagens") || url.includes("/api/inventory-counts")) {
    if (Array.isArray(data)) {
      return data.map((item) => ({
        id: item.id,
        itemCode: item.codigo_item,
        itemDescription: item.descricao_item,
        location: item.local,
        quantity: item.quantidade.toString(),
        volumeCount: item.volumes,
        timestamp: item.data_contagem,
      }));
    }
  }

  // Para outros dados ou se já estão no formato correto
  return data;
}

export const getQueryFn: <T>(options: {
  on401: UnauthorizedBehavior;
}) => QueryFunction<T> =
  ({ on401: unauthorizedBehavior }) =>
    async ({ queryKey }) => {
      const url = queryKey.join("/") as string;
      // Construir URL completa para a API Go
      const fullUrl = url.startsWith("http") ? url : `${API_BASE_URL}${url}`;

      const res = await fetch(fullUrl, {
        // credentials: "include",
      });

      if (unauthorizedBehavior === "returnNull" && res.status === 401) {
        return null;
      }

      await throwIfResNotOk(res);
      const jsonData = await res.json();

      // Se a resposta tem estrutura {success, message, data}, extrair apenas data
      if (
        jsonData &&
        typeof jsonData === "object" &&
        "success" in jsonData &&
        "data" in jsonData
      ) {
        const processedData = processApiResponse(jsonData.data, url);
        return processedData;
      }

      const processedData = processApiResponse(jsonData, url);
      return processedData;
    };

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      queryFn: getQueryFn({ on401: "throw" }),
      refetchInterval: false,
      refetchOnWindowFocus: false,
      staleTime: Infinity,
      retry: false,
    },
    mutations: {
      retry: false,
    },
  },
});
