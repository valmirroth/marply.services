import { type Item, type InsertItem, type InventoryCount, type InsertInventoryCount, type UpdateInventoryCount } from "@shared/schema";
import { randomUUID } from "crypto";

export interface IStorage {
  // Item operations
  getItemByCode(code: string): Promise<Item | undefined>;
  createItem(item: InsertItem): Promise<Item>;
  
  // Inventory count operations
  getInventoryCounts(): Promise<InventoryCount[]>;
  getInventoryCount(id: string): Promise<InventoryCount | undefined>;
  createInventoryCount(count: InsertInventoryCount): Promise<InventoryCount>;
  updateInventoryCount(update: UpdateInventoryCount): Promise<InventoryCount | undefined>;
  deleteInventoryCount(id: string): Promise<boolean>;
}

export class MemStorage implements IStorage {
  private items: Map<string, Item>;
  private inventoryCounts: Map<string, InventoryCount>;

  constructor() {
    this.items = new Map();
    this.inventoryCounts = new Map();
    
    // Add some sample items for demonstration
    this.seedItems();
  }

  private seedItems() {
    const sampleItems = [
      { code: "A001234", description: "Parafuso Sextavado M12x50 Aço Inox" },
      { code: "B002456", description: "Arruela Lisa M12 Aço Galvanizado" },
      { code: "C003789", description: "Porca Sextavada M12 Aço Inox" },
      { code: "D004567", description: "Parafuso Phillips M8x40 Aço Carbono" },
      { code: "E005123", description: "Porca Flangeada M10 Aço Inox" }
    ];

    sampleItems.forEach(item => {
      const id = randomUUID();
      const itemWithId: Item = { id, ...item };
      this.items.set(item.code, itemWithId);
    });
  }

  async getItemByCode(code: string): Promise<Item | undefined> {
    return this.items.get(code);
  }

  async createItem(insertItem: InsertItem): Promise<Item> {
    const id = randomUUID();
    const item: Item = { id, ...insertItem };
    this.items.set(insertItem.code, item);
    return item;
  }

  async getInventoryCounts(): Promise<InventoryCount[]> {
    return Array.from(this.inventoryCounts.values())
      .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime());
  }

  async getInventoryCount(id: string): Promise<InventoryCount | undefined> {
    return this.inventoryCounts.get(id);
  }

  async createInventoryCount(insertCount: InsertInventoryCount): Promise<InventoryCount> {
    const id = randomUUID();
    const count: InventoryCount = {
      id,
      ...insertCount,
      timestamp: new Date(),
    };
    this.inventoryCounts.set(id, count);
    return count;
  }

  async updateInventoryCount(update: UpdateInventoryCount): Promise<InventoryCount | undefined> {
    const existing = this.inventoryCounts.get(update.id);
    if (!existing) return undefined;

    const updated: InventoryCount = {
      ...existing,
      ...update,
      timestamp: new Date(),
    };
    this.inventoryCounts.set(update.id, updated);
    return updated;
  }

  async deleteInventoryCount(id: string): Promise<boolean> {
    return this.inventoryCounts.delete(id);
  }
}

export const storage = new MemStorage();
