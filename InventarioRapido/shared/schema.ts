import { sql } from "drizzle-orm";
import { pgTable, text, varchar, decimal, integer, timestamp } from "drizzle-orm/pg-core";
import { createInsertSchema } from "drizzle-zod";
import { z } from "zod";

export const items = pgTable("items", {
  id: varchar("id").primaryKey().default(sql`gen_random_uuid()`),
  code: text("code").notNull().unique(),
  description: text("description").notNull(),
});

export const inventoryCounts = pgTable("inventory_counts", {
  id: varchar("id").primaryKey().default(sql`gen_random_uuid()`),
  itemCode: text("item_code").notNull(),
  itemDescription: text("item_description").notNull(),
  location: text("location").notNull(),
  quantity: decimal("quantity", { precision: 10, scale: 5 }).notNull(),
  volumeCount: decimal("volume_count", { precision: 10, scale: 5 }).notNull(),
  empresa: integer("empresa").notNull(),
  timestamp: timestamp("timestamp").defaultNow().notNull(),
});

export const insertItemSchema = createInsertSchema(items).omit({
  id: true,
});

export const insertInventoryCountSchema = createInsertSchema(inventoryCounts).omit({
  id: true,
  timestamp: true,
});

export const updateInventoryCountSchema = insertInventoryCountSchema.partial().extend({
  id: z.string(),
});

export type InsertItem = z.infer<typeof insertItemSchema>;
export type Item = typeof items.$inferSelect;

export type InsertInventoryCount = z.infer<typeof insertInventoryCountSchema>;
export type UpdateInventoryCount = z.infer<typeof updateInventoryCountSchema>;
export type InventoryCount = typeof inventoryCounts.$inferSelect;
