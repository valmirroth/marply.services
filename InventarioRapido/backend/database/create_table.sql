-- Criação da tabela CST_CONTAGEM_LOCAL
-- Execute este script no SQL Server antes de iniciar a API

USE [InventoryDB];
GO

-- Criar a tabela se não existir
IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='CST_CONTAGEM_LOCAL' AND xtype='U')
BEGIN
    CREATE TABLE [dbo].[CST_CONTAGEM_LOCAL] (
        [ID] [uniqueidentifier] NOT NULL DEFAULT NEWID(),
        [CODIGO_ITEM] [varchar](50) NOT NULL,
        [DESCRICAO_ITEM] [varchar](255) NOT NULL,
        [LOCAL] [varchar](100) NOT NULL,
        [QUANTIDADE] [decimal](10,5) NOT NULL,
        [VOLUMES] [int] NOT NULL,
        [DATA_CONTAGEM] [datetime2](7) NOT NULL DEFAULT GETDATE(),
        [USUARIO_CONTAGEM] [varchar](100) NULL DEFAULT 'Sistema',
        [STATUS] [varchar](20) NOT NULL DEFAULT 'ATIVO',
        CONSTRAINT [PK_CST_CONTAGEM_LOCAL] PRIMARY KEY CLUSTERED ([ID] ASC)
    );
    
    -- Criar índices para melhor performance
    CREATE INDEX [IX_CST_CONTAGEM_LOCAL_CODIGO_ITEM] ON [dbo].[CST_CONTAGEM_LOCAL] ([CODIGO_ITEM]);
    CREATE INDEX [IX_CST_CONTAGEM_LOCAL_LOCAL] ON [dbo].[CST_CONTAGEM_LOCAL] ([LOCAL]);
    CREATE INDEX [IX_CST_CONTAGEM_LOCAL_DATA] ON [dbo].[CST_CONTAGEM_LOCAL] ([DATA_CONTAGEM]);
    
    PRINT 'Tabela CST_CONTAGEM_LOCAL criada com sucesso!';
END
ELSE
BEGIN
    PRINT 'Tabela CST_CONTAGEM_LOCAL já existe.';
END
GO