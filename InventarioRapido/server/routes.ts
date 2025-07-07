import type { Express } from "express";
import { createServer, type Server } from "http";
import { storage } from "./storage";
import { insertInventoryCountSchema, updateInventoryCountSchema } from "@shared/schema";

export async function registerRoutes(app: Express): Promise<Server> {
  // Get item by code
  app.get("/api/items/:code", async (req, res) => {
    try {
      const { code } = req.params;
      const item = await storage.getItemByCode(code);
      
      if (!item) {
        return res.status(404).json({ message: "Item não encontrado" });
      }
      
      res.json(item);
    } catch (error) {
      res.status(500).json({ message: "Erro interno do servidor" });
    }
  });

  // Get all inventory counts
  app.get("/api/inventory-counts", async (req, res) => {
    try {
      const counts = await storage.getInventoryCounts();
      res.json(counts);
    } catch (error) {
      res.status(500).json({ message: "Erro interno do servidor" });
    }
  });

  // Create new inventory count
  app.post("/api/inventory-counts", async (req, res) => {
    try {
      const result = insertInventoryCountSchema.safeParse(req.body);
      
      if (!result.success) {
        return res.status(400).json({ 
          message: "Dados inválidos",
          errors: result.error.errors 
        });
      }

      const count = await storage.createInventoryCount(result.data);
      res.status(201).json(count);
    } catch (error) {
      res.status(500).json({ message: "Erro interno do servidor" });
    }
  });

  // Update inventory count
  app.put("/api/inventory-counts/:id", async (req, res) => {
    try {
      const { id } = req.params;
      const result = updateInventoryCountSchema.safeParse({ ...req.body, id });
      
      if (!result.success) {
        return res.status(400).json({ 
          message: "Dados inválidos",
          errors: result.error.errors 
        });
      }

      const count = await storage.updateInventoryCount(result.data);
      
      if (!count) {
        return res.status(404).json({ message: "Contagem não encontrada" });
      }
      
      res.json(count);
    } catch (error) {
      res.status(500).json({ message: "Erro interno do servidor" });
    }
  });

  // Delete inventory count
  app.delete("/api/inventory-counts/:id", async (req, res) => {
    try {
      const { id } = req.params;
      const deleted = await storage.deleteInventoryCount(id);
      
      if (!deleted) {
        return res.status(404).json({ message: "Contagem não encontrada" });
      }
      
      res.json({ message: "Contagem excluída com sucesso" });
    } catch (error) {
      res.status(500).json({ message: "Erro interno do servidor" });
    }
  });

  const httpServer = createServer(app);
  return httpServer;
}
