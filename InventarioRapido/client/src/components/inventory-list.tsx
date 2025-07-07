import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import type { InventoryCount } from "@shared/schema";

interface InventoryListProps {
  counts: InventoryCount[];
  onEdit: (count: InventoryCount) => void;
  onDelete: (count: InventoryCount) => void;
  disabled?: boolean;
}

export default function InventoryList({ counts, onEdit, onDelete, disabled = false }: InventoryListProps) {
  // Garantir que counts seja sempre um array
  const safeCount = Array.isArray(counts) ? counts : [];
  
  if (safeCount.length === 0) {
    return (
      <Card className="shadow-lg">
        <div className="bg-gray-50 border-b px-8 py-6">
          <h2 className="text-xl font-semibold text-gray-900 flex items-center">
            <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Últimas Contagens
            <span className="ml-2 bg-blue-600 text-white text-sm px-2 py-1 rounded-full">0</span>
          </h2>
        </div>
        <CardContent className="p-12 text-center">
          <div className="text-gray-400 mb-6">
            <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
          </div>
          <h3 className="text-xl font-semibold text-gray-600 mb-2">Nenhuma contagem registrada</h3>
          <p className="text-gray-500">Comece registrando sua primeira contagem de inventário.</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="shadow-lg overflow-hidden">
      <div className="bg-gray-50 border-b px-8 py-6">
        <h2 className="text-xl font-semibold text-gray-900 flex items-center">
          <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Últimas Contagens
          <span className="ml-2 bg-blue-600 text-white text-sm px-2 py-1 rounded-full">{safeCount.length}</span>
        </h2>
      </div>

      <div className="divide-y divide-gray-200">
        {safeCount.map((count) => (
          <div key={count.id} className="p-8 hover:bg-gray-50 transition-colors" data-testid={`count-item-${count.id}`}>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
              <div className="text-center">
                <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">Código</div>
                <div className="text-lg font-semibold text-gray-900" data-testid={`text-code-${count.id}`}>{count.itemCode}</div>
              </div>
              <div className="text-center">
                <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">Local</div>
                <div className="text-lg font-semibold text-gray-900" data-testid={`text-location-${count.id}`}>{count.location}</div>
              </div>
              <div className="text-center">
                <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">Quantidade</div>
                <div className="text-lg font-semibold text-gray-900" data-testid={`text-quantity-${count.id}`}>{count.quantity}</div>
              </div>
              <div className="text-center">
                <div className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-1">Volumes</div>
                <div className="text-lg font-semibold text-gray-900" data-testid={`text-volumes-${count.id}`}>{count.volumeCount}</div>
              </div>
            </div>

            <div className="mb-4">
              <strong>Descrição: </strong>
              <span data-testid={`text-description-${count.id}`}>{count.itemDescription}</span>
            </div>

            <div className="flex justify-between items-center">
              <div className="text-sm text-gray-500 flex items-center">
                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                {new Date(count.timestamp).toLocaleString('pt-BR')}
              </div>

              {!disabled && (
                <div className="flex gap-2">
                  <Button
                    variant="default"
                    size="sm"
                    onClick={() => onEdit(count)}
                    data-testid={`button-edit-${count.id}`}
                    className="bg-green-600 hover:bg-green-700"
                  >
                    <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    Editar
                  </Button>
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => onDelete(count)}
                    data-testid={`button-delete-${count.id}`}
                  >
                    <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Excluir
                  </Button>
                </div>
              )}
              {disabled && (
                <div className="text-sm text-gray-500 bg-gray-100 px-3 py-2 rounded-md">
                  <svg className="w-4 h-4 mr-1 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                  Bloqueado
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </Card>
  );
}
