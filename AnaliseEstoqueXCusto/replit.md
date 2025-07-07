# Analisador de Custo Médio - SQL Server

## Overview

This is a standalone Python script designed for **local execution on user's server environment**. It calculates average cost of items based on SQL Server movement data, compares values with original records, and sends email notifications when significant differences are detected. The system is designed as a monitoring and auditing tool for inventory cost tracking in business operations.

**Important**: This is NOT a cloud-based or Replit-hosted application. The code is provided here for download and deployment to the user's local infrastructure where it will run as a scheduled task (cron/Task Scheduler) daily at 5:00 AM.

## User Preferences

Preferred communication style: Simple, everyday language.

## System Architecture

### Application Type
**Standalone Python Script for Local Deployment** - This is a command-line application designed to run as a scheduled task or cron job on the user's local server, not a web application or cloud service. It performs batch processing of data and sends email reports.

**Deployment Model**: 
- User downloads the code from this Replit
- Installs it on their local Windows/Linux server
- Configures credentials in `.env` file
- Schedules execution via Windows Task Scheduler or Linux cron
- Script runs daily at 5:00 AM in the user's timezone

### Core Components

1. **Data Processing Layer**
   - **Technology**: Python 3.8+ with pandas for data manipulation
   - **Purpose**: Fetches movement data from SQL Server, calculates average costs, and compares with stored values
   - **Key Class**: `CustoMedioAnalyzer` - main orchestrator for all operations
   - **Design Pattern**: Single responsibility class that handles database connection, data processing, and email notification

2. **Database Integration**
   - **Technology**: SQL Server with pyodbc driver
   - **Connection Method**: ODBC Driver 17 for SQL Server
   - **Access Pattern**: Read-only queries to fetch movement data and compare with existing cost values
   - **Configuration**: Connection details stored in environment variables (server, database, credentials, port)

3. **Notification System**
   - **Technology**: SMTP email using Python's built-in email libraries
   - **Purpose**: Sends automated reports when cost differences exceed configured threshold
   - **Configuration**: SMTP server details and recipient lists managed via environment variables
   - **Threshold Logic**: Configurable difference threshold (default 0.05 or 5%) to filter significant discrepancies

4. **Configuration Management**
   - **Technology**: python-dotenv for environment variable management
   - **Security**: Sensitive credentials stored in `.env` file (not committed to repository)
   - **Pattern**: All configuration externalized from code for easy deployment across environments

### Data Flow

1. Script loads environment variables from `.env` file
2. Establishes connection to SQL Server using ODBC driver
3. Queries movement data from database
4. Calculates average costs using pandas DataFrame operations
5. Compares calculated values against stored database values
6. Filters differences exceeding threshold
7. Generates email report with findings
8. Sends notification via SMTP to configured recipients

### Error Handling Strategy
- Connection validation before main processing
- Test script (`test_conexao.py`) provided for pre-execution validation
- Environment variable verification to ensure all required config is present

## External Dependencies

### Database Systems
- **SQL Server** - Primary data source for movement and cost data
  - Requires ODBC Driver 17 for SQL Server
  - Default port: 1433 (configurable)
  - Read-only access pattern

### Email Services
- **SMTP Server** - Email notification delivery
  - Default port: 587 (configurable)
  - Supports authenticated SMTP

### Python Libraries
- **pyodbc 5.0.1** - SQL Server database connectivity via ODBC
- **pandas 2.1.4** - Data manipulation and analysis
- **python-dotenv 1.0.0** - Environment variable management

### System Dependencies
- **ODBC Driver 17 for SQL Server** - Required system-level driver for database connectivity
  - Windows: Microsoft installer available
  - Linux: unixodbc-dev package required

### Environment Variables Required
- `SQL_SERVER` - Database server hostname/IP
- `SQL_DATABASE` - Database name
- `SQL_USERNAME` - Database user credentials
- `SQL_PASSWORD` - Database password
- `SQL_PORT` - Database port (optional, defaults to 1433)
- `SMTP_SERVER` - Email server hostname
- `SMTP_PORT` - Email server port (optional, defaults to 587)
- `SMTP_USERNAME` - Email authentication username
- `SMTP_PASSWORD` - Email authentication password
- `EMAIL_FROM` - Sender email address
- `EMAIL_TO` - Recipient email address (defaults to ti01@marply.com.br)
- `DIFERENCA_THRESHOLD` - Cost difference threshold for reporting (optional, defaults to 0.05)