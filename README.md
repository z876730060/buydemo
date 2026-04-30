# Buy-Demo ERP 系统

> **项目定位**：轻量级ERP系统，覆盖采购→库存→销售→财务全链路
> **技术栈**：Go + Gin + GORM + SQLite + Vue.js(CDN) + Element UI + JWT + ECharts
> **架构风格**：前后端分离（SPA），RESTful API
> **版本**：v2.1 | **更新**：2026-05-01

---

## 快速开始

### 本地运行
```bash
git clone https://github.com/z876730060/buydemo.git
cd buydemo
go build -o buydemo.exe .
./buydemo.exe
# 访问 http://localhost:8080
# 默认账号: admin / admin123
```

### Docker 部署
```bash
docker-compose up -d
# 访问 http://localhost:8080
# 数据持久化在 ./data/ 目录
```

---

## 模块功能

| 模块 | 功能 | 状态 |
|:---|:---|:---:|
| **认证** | 登录/登出/JWT签发/修改密码 | ✅ |
| **供应商** | CRUD/分页搜索/采购单关联/应付账款汇总 | ✅ |
| **客户** | CRUD/分页搜索/销售单关联/应收账款汇总 | ✅ |
| **商品** | CRUD/分类/分页搜索/库存关联/采购销售流水 | ✅ |
| **采购单** | 创建/编辑/审核/入库/取消/加权平均成本 | ✅ |
| **销售单** | 创建/编辑/审核/出库/取消/自动扣库存 | ✅ |
| **库存** | 台账/流水/低库存预警/手动调整 | ✅ |
| **财务管理** | 应付/应收/费用/收付款记录/收支趋势 | ✅ |
| **用户管理** | CRUD/密码重置/角色权限 | ✅ |
| **操作日志** | 全量自动记录/筛选/CSV导出 | ✅ |
| **仪表盘** | ECharts图表/状态分布/应收应付汇总 | ✅ |
| **报表中心** | 采购/销售/库存报表 | ✅ |
| **数据导出** | 供应商/商品/订单/库存CSV导出 | ✅ |
| **数据可视化** | ECharts饼图/柱图/趋势线图 | ✅ |
| **数据导入** | CSV批量导入商品/供应商 | ✅ |
| **打印功能** | 采购单/销售单打印 | ✅ |

---

## 技术架构

```
buy-demo/
├── main.go
├── config/config.go              # 配置管理
├── database/database.go          # 数据库初始化 & 迁移
├── models/                       # 数据模型
│   ├── user.go, supplier.go, customer.go
│   ├── product.go, purchase_order.go
│   ├── sales.go, inventory.go
│   ├── finance.go, operation_log.go
├── handlers/                     # API处理器
│   ├── auth.go, user.go
│   ├── supplier.go, customer.go, product.go
│   ├── purchase.go, sales.go, inventory.go
│   ├── finance.go, operation_log.go, report.go
├── middlewares/
│   ├── auth.go                   # JWT认证
│   └── operation_log.go          # 操作日志
├── router/router.go              # 路由注册
├── static/index.html             # SPA前端
└── README.md
```

---

## API 总览（50+ 接口）

### 认证
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| POST | /api/auth/login | 登录 |
| GET | /api/auth/me | 当前用户 |
| POST | /api/auth/change-password | 修改密码 |

### 供应商
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| GET | /api/suppliers | 列表（分页+搜索） |
| GET | /api/suppliers/all | 全部启用 |
| GET | /api/suppliers/:id | 详情 |
| GET | /api/suppliers/:id/orders | 采购单+应付汇总 |
| POST | /api/suppliers | 新增 |
| PUT | /api/suppliers/:id | 编辑 |
| DELETE | /api/suppliers/:id | 删除 |
| POST | /api/suppliers/import | CSV导入 |

### 客户
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| GET | /api/customers | 列表（分页+搜索） |
| GET | /api/customers/all | 全部启用 |
| GET | /api/customers/:id | 详情 |
| GET | /api/customers/:id/orders | 销售单+应收汇总 |
| POST/PUT/DELETE | /api/customers/* | CRUD |

### 商品
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| GET | /api/products | 列表（分页+搜索+分类） |
| GET | /api/products/all | 全部启用 |
| GET | /api/products/:id | 详情 |
| GET | /api/products/:id/detail | 完整信息（库存+采购+销售+流水） |
| POST | /api/products | 新增（自动创建库存） |
| PUT | /api/products/:id | 编辑 |
| DELETE | /api/products/:id | 删除 |
| POST | /api/products/import | CSV导入 |

### 采购单
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| GET | /api/purchase-orders | 列表 |
| GET | /api/purchase-orders/:id | 详情 |
| POST | /api/purchase-orders | 创建 |
| PUT | /api/purchase-orders/:id | 编辑（草稿） |
| POST | /api/purchase-orders/:id/approve | 审核 |
| POST | /api/purchase-orders/:id/receive | 入库（更新成本价） |
| POST | /api/purchase-orders/:id/cancel | 取消 |

### 销售单
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| GET | /api/sales-orders | 列表 |
| GET | /api/sales-orders/:id | 详情 |
| POST | /api/sales-orders | 创建 |
| PUT | /api/sales-orders/:id | 编辑（草稿） |
| POST | /api/sales-orders/:id/approve | 审核 |
| POST | /api/sales-orders/:id/deliver | 出库（扣库存） |
| POST | /api/sales-orders/:id/cancel | 取消 |

### 库存
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| GET | /api/inventories | 列表 |
| GET | /api/inventories/logs | 流水 |
| GET | /api/inventories/low-stock | 低库存预警 |
| POST | /api/inventories/adjust | 手动调整 |

### 财务管理
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| GET | /api/finance/summary | 概览 |
| GET | /api/finance/payable | 应付列表 |
| POST | /api/finance/payable/:id/pay | 付款 |
| GET | /api/finance/receivable | 应收列表 |
| POST | /api/finance/receivable/:id/receive | 收款 |
| CRUD | /api/finance/expenses | 费用管理 |
| GET | /api/finance/payments | 收付款记录 |

### 其他
| 方法 | 路径 | 说明 |
|:---|:---|:---|
| CRUD | /api/users | 用户管理 |
| GET | /api/operation-logs | 操作日志（筛选+导出） |
| GET | /api/reports/* | 采购/销售/库存报表 |
| GET | /api/export/* | CSV导出 |
| GET | /api/dashboard/stats | 仪表盘数据 |

---

## 核心业务流程

```
采购到付款: 创建采购单 → 审核 → 入库(加权平均成本) → 自动应付 → 分期付款
销售到收款: 创建销售单 → 审核 → 出库(扣库存) → 自动应收 → 分期收款
库存变动:   采购入库(+) → 销售出库(-) → 手动调整(±) → 完整流水追溯
数据关联:   商品↔库存↔采购↔供应商↔应付↔付款
            商品↔库存↔销售↔客户↔应收↔收款
```

---

## 数据库模型

| 表 | 说明 |
|:---|:---|
| users | 用户（admin/operator） |
| suppliers | 供应商（编码唯一） |
| customers | 客户（编码唯一） |
| products | 商品/物料（编码唯一，自动创建库存） |
| purchase_orders/items | 采购单及明细 |
| sales_orders/items | 销售单及明细 |
| inventories | 库存台账（与商品1:1） |
| inventory_logs | 库存流水（in/out/adjust） |
| account_payables | 应付账款（采购入库自动生成） |
| account_receivables | 应收账款（销售出库自动生成） |
| expenses | 费用记录 |
| payment_records | 收付款流水 |
| operation_logs | 操作日志（全量自动记录） |

---

## 开发计划

### ✅ 第一期：基础框架
1. ✅ 项目初始化 & 目录结构搭建
2. ✅ 数据库初始化 & 模型定义
3. ✅ 认证模块（登录/JWT）
4. ✅ 供应商、商品CRUD
5. ✅ 采购单流转（创建→审核→入库）
6. ✅ 库存台账管理
7. ✅ 前端SPA界面

### ✅ 第二期：业务扩展
8. ✅ 销售管理（客户、销售订单、出库）
9. ✅ 财务管理（应收/应付/费用/收付款）
10. ✅ 用户管理 & 操作日志
11. ✅ 报表中心 & 数据导出

### ✅ 第三期：全链路集成
12. ✅ 商品详情（库存+采购+销售+流水）
13. ✅ 供应商/客户详情（订单+财务汇总）
14. ✅ 跨表点击跳转
15. ✅ 采购入库加权平均成本
16. ✅ 库存手动调整
17. ✅ ECharts可视化仪表盘
18. ✅ CSV批量导入（商品/供应商）
19. ✅ 采购单/销售单打印

### ✅ 第四期：部署与完善
20. ✅ Docker部署（Dockerfile + docker-compose）
21. 🔜 数据库备份与恢复
22. 🔜 系统参数配置
23. 🔜 API文档（Swagger/OpenAPI）
