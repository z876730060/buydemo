package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "url": "https://github.com/z876730060/buydemo"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/login": {
            "post": {
                "description": "用户登录系统",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["认证"],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "登录信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "username": {"type": "string", "example": "admin"},
                                "password": {"type": "string", "example": "admin123"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "token": {"type": "string"},
                                "user": {
                                    "type": "object",
                                    "properties": {
                                        "id": {"type": "integer"},
                                        "username": {"type": "string"},
                                        "real_name": {"type": "string"},
                                        "role": {"type": "string"}
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/auth/me": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取当前登录用户信息",
                "produces": ["application/json"],
                "tags": ["认证"],
                "summary": "获取当前用户",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "user": {
                                    "type": "object",
                                    "properties": {
                                        "id": {"type": "integer"},
                                        "username": {"type": "string"},
                                        "real_name": {"type": "string"},
                                        "role": {"type": "string"}
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/auth/change-password": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "修改当前用户密码",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["认证"],
                "summary": "修改密码",
                "parameters": [
                    {
                        "description": "密码信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "old_password": {"type": "string"},
                                "new_password": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {"description": "修改成功"}
                }
            }
        },
        "/api/suppliers": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取供应商列表，支持分页和搜索",
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "供应商列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "keyword", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/Supplier"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建新供应商",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "新增供应商",
                "parameters": [
                    {"description": "供应商信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Supplier"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/suppliers/all": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取所有启用状态的供应商",
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "全部供应商",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/Supplier"}
                        }
                    }
                }
            }
        },
        "/api/suppliers/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取供应商详情",
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "供应商详情",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/Supplier"}
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新供应商信息",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "编辑供应商",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "供应商信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Supplier"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "删除供应商",
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "删除供应商",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "删除成功"}}
            }
        },
        "/api/suppliers/{id}/orders": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取供应商的采购单和应付账款汇总",
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "供应商订单汇总",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "orders": {"type": "array", "items": {"$ref": "#/definitions/PurchaseOrder"}},
                                "total_payable": {"type": "number"}
                            }
                        }
                    }
                }
            }
        },
        "/api/suppliers/import": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "通过CSV文件批量导入供应商",
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "tags": ["供应商"],
                "summary": "导入供应商",
                "parameters": [
                    {"name": "file", "in": "formData", "required": true, "type": "file", "description": "CSV文件"}
                ],
                "responses": {"200": {"description": "导入成功"}}
            }
        },
        "/api/customers": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取客户列表，支持分页和搜索",
                "produces": ["application/json"],
                "tags": ["客户"],
                "summary": "客户列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "keyword", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/Customer"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建新客户",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["客户"],
                "summary": "新增客户",
                "parameters": [
                    {"description": "客户信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Customer"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/customers/all": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取所有启用状态的客户",
                "produces": ["application/json"],
                "tags": ["客户"],
                "summary": "全部客户",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/Customer"}
                        }
                    }
                }
            }
        },
        "/api/customers/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取客户详情",
                "produces": ["application/json"],
                "tags": ["客户"],
                "summary": "客户详情",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/Customer"}
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新客户信息",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["客户"],
                "summary": "编辑客户",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "客户信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Customer"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "删除客户",
                "produces": ["application/json"],
                "tags": ["客户"],
                "summary": "删除客户",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "删除成功"}}
            }
        },
        "/api/customers/{id}/orders": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取客户的销售单和应收账款汇总",
                "produces": ["application/json"],
                "tags": ["客户"],
                "summary": "客户订单汇总",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "orders": {"type": "array", "items": {"$ref": "#/definitions/SalesOrder"}},
                                "total_receivable": {"type": "number"}
                            }
                        }
                    }
                }
            }
        },
        "/api/products": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取商品列表，支持分页、搜索和分类筛选",
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "商品列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "keyword", "in": "query", "type": "string"},
                    {"name": "category", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/Product"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建新商品，自动创建库存记录",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "新增商品",
                "parameters": [
                    {"description": "商品信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Product"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/products/all": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取所有启用状态的商品",
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "全部商品",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/Product"}
                        }
                    }
                }
            }
        },
        "/api/products/categories": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取商品分类列表",
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "商品分类",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"type": "string"}
                        }
                    }
                }
            }
        },
        "/api/products/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取商品详情",
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "商品详情",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/Product"}
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新商品信息",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "编辑商品",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "商品信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Product"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "删除商品",
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "删除商品",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "删除成功"}}
            }
        },
        "/api/products/{id}/detail": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取商品完整信息，包含库存、采购、销售、流水",
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "商品完整信息",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "product": {"$ref": "#/definitions/Product"},
                                "inventory": {"$ref": "#/definitions/Inventory"},
                                "purchase_orders": {"type": "array", "items": {"$ref": "#/definitions/PurchaseOrder"}},
                                "sales_orders": {"type": "array", "items": {"$ref": "#/definitions/SalesOrder"}},
                                "inventory_logs": {"type": "array", "items": {"$ref": "#/definitions/InventoryLog"}}
                            }
                        }
                    }
                }
            }
        },
        "/api/products/import": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "通过CSV文件批量导入商品",
                "consumes": ["multipart/form-data"],
                "produces": ["application/json"],
                "tags": ["商品"],
                "summary": "导入商品",
                "parameters": [
                    {"name": "file", "in": "formData", "required": true, "type": "file", "description": "CSV文件"}
                ],
                "responses": {"200": {"description": "导入成功"}}
            }
        },
        "/api/purchase-orders": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取采购单列表",
                "produces": ["application/json"],
                "tags": ["采购单"],
                "summary": "采购单列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "status", "in": "query", "type": "string"},
                    {"name": "supplier_id", "in": "query", "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/PurchaseOrder"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建采购单",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["采购单"],
                "summary": "新建采购单",
                "parameters": [
                    {"description": "采购单信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/PurchaseOrder"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/purchase-orders/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取采购单详情",
                "produces": ["application/json"],
                "tags": ["采购单"],
                "summary": "采购单详情",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/PurchaseOrder"}
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新采购单（仅草稿状态）",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["采购单"],
                "summary": "编辑采购单",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "采购单信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/PurchaseOrder"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            }
        },
        "/api/purchase-orders/{id}/approve": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "审核采购单",
                "produces": ["application/json"],
                "tags": ["采购单"],
                "summary": "审核采购单",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "审核成功"}}
            }
        },
        "/api/purchase-orders/{id}/receive": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "采购入库，更新库存和成本价（加权平均）",
                "produces": ["application/json"],
                "tags": ["采购单"],
                "summary": "采购入库",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "入库成功"}}
            }
        },
        "/api/purchase-orders/{id}/cancel": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "取消采购单",
                "produces": ["application/json"],
                "tags": ["采购单"],
                "summary": "取消采购单",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "取消成功"}}
            }
        },
        "/api/sales-orders": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取销售单列表",
                "produces": ["application/json"],
                "tags": ["销售单"],
                "summary": "销售单列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "status", "in": "query", "type": "string"},
                    {"name": "customer_id", "in": "query", "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/SalesOrder"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建销售单",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["销售单"],
                "summary": "新建销售单",
                "parameters": [
                    {"description": "销售单信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/SalesOrder"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/sales-orders/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取销售单详情",
                "produces": ["application/json"],
                "tags": ["销售单"],
                "summary": "销售单详情",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/SalesOrder"}
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新销售单（仅草稿状态）",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["销售单"],
                "summary": "编辑销售单",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "销售单信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/SalesOrder"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            }
        },
        "/api/sales-orders/{id}/approve": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "审核销售单",
                "produces": ["application/json"],
                "tags": ["销售单"],
                "summary": "审核销售单",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "审核成功"}}
            }
        },
        "/api/sales-orders/{id}/deliver": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "销售出库，扣减库存",
                "produces": ["application/json"],
                "tags": ["销售单"],
                "summary": "销售出库",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "出库成功"}}
            }
        },
        "/api/sales-orders/{id}/cancel": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "取消销售单",
                "produces": ["application/json"],
                "tags": ["销售单"],
                "summary": "取消销售单",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "取消成功"}}
            }
        },
        "/api/inventories": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取库存台账列表",
                "produces": ["application/json"],
                "tags": ["库存"],
                "summary": "库存列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "keyword", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/Inventory"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            }
        },
        "/api/inventories/logs": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取库存流水记录",
                "produces": ["application/json"],
                "tags": ["库存"],
                "summary": "库存流水",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "product_id", "in": "query", "type": "integer"},
                    {"name": "type", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/InventoryLog"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            }
        },
        "/api/inventories/low-stock": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取低库存预警商品",
                "produces": ["application/json"],
                "tags": ["库存"],
                "summary": "低库存预警",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/Inventory"}
                        }
                    }
                }
            }
        },
        "/api/inventories/adjust": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "手动调整库存",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["库存"],
                "summary": "库存调整",
                "parameters": [
                    {
                        "description": "调整信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "product_id": {"type": "integer"},
                                "quantity": {"type": "integer"},
                                "reason": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {"200": {"description": "调整成功"}}
            }
        },
        "/api/finance/summary": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取财务概览",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "财务概览",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "total_payable": {"type": "number"},
                                "total_receivable": {"type": "number"},
                                "total_expenses": {"type": "number"},
                                "recent_payments": {"type": "array", "items": {"$ref": "#/definitions/PaymentRecord"}}
                            }
                        }
                    }
                }
            }
        },
        "/api/finance/payable": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取应付账款列表",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "应付列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/AccountPayable"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            }
        },
        "/api/finance/payable/{id}/pay": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "付款",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "付款",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {
                        "description": "付款信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "amount": {"type": "number"},
                                "payment_method": {"type": "string"},
                                "remark": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {"200": {"description": "付款成功"}}
            }
        },
        "/api/finance/receivable": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取应收账款列表",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "应收列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/AccountReceivable"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            }
        },
        "/api/finance/receivable/{id}/receive": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "收款",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "收款",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {
                        "description": "收款信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "amount": {"type": "number"},
                                "payment_method": {"type": "string"},
                                "remark": {"type": "string"}
                            }
                        }
                    }
                ],
                "responses": {"200": {"description": "收款成功"}}
            }
        },
        "/api/finance/expenses": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取费用记录列表",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "费用列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/Expense"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建费用记录",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "新增费用",
                "parameters": [
                    {"description": "费用信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Expense"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/finance/expenses/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新费用记录",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "编辑费用",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "费用信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Expense"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "删除费用记录",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "删除费用",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "删除成功"}}
            }
        },
        "/api/finance/payments": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取收付款记录",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "收付款记录",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "type", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/PaymentRecord"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            }
        },
        "/api/finance/payment-methods": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取付款方式列表",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "付款方式",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"type": "string"}
                        }
                    }
                }
            }
        },
        "/api/finance/expense-categories": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取费用类别列表",
                "produces": ["application/json"],
                "tags": ["财务管理"],
                "summary": "费用类别",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"type": "string"}
                        }
                    }
                }
            }
        },
        "/api/users": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取用户列表（仅管理员）",
                "produces": ["application/json"],
                "tags": ["用户管理"],
                "summary": "用户列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/User"}
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建新用户（仅管理员）",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["用户管理"],
                "summary": "新增用户",
                "parameters": [
                    {"description": "用户信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/User"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/users/{id}": {
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新用户信息（仅管理员）",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["用户管理"],
                "summary": "编辑用户",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "用户信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/User"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "删除用户（仅管理员）",
                "produces": ["application/json"],
                "tags": ["用户管理"],
                "summary": "删除用户",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "删除成功"}}
            }
        },
        "/api/users/{id}/reset-password": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "重置用户密码（仅管理员）",
                "produces": ["application/json"],
                "tags": ["用户管理"],
                "summary": "重置密码",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "重置成功"}}
            }
        },
        "/api/operation-logs": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取操作日志列表",
                "produces": ["application/json"],
                "tags": ["操作日志"],
                "summary": "操作日志",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20},
                    {"name": "action", "in": "query", "type": "string"},
                    {"name": "user_id", "in": "query", "type": "integer"},
                    {"name": "start_date", "in": "query", "type": "string"},
                    {"name": "end_date", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/OperationLog"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            }
        },
        "/api/operation-logs/export": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "导出操作日志为CSV",
                "produces": ["text/csv"],
                "tags": ["操作日志"],
                "summary": "导出日志",
                "responses": {
                    "200": {
                        "description": "CSV文件下载"
                    }
                }
            }
        },
        "/api/dashboard/stats": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取仪表盘统计数据",
                "produces": ["application/json"],
                "tags": ["仪表盘"],
                "summary": "仪表盘数据",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "customers": {"type": "integer"},
                                "suppliers": {"type": "integer"},
                                "products": {"type": "integer"},
                                "today_orders": {"type": "integer"},
                                "today_sales": {"type": "number"}
                            }
                        }
                    }
                }
            }
        },
        "/api/reports/purchase": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取采购报表",
                "produces": ["application/json"],
                "tags": ["报表"],
                "summary": "采购报表",
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string"},
                    {"name": "end_date", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "start_date": {"type": "string"},
                                "end_date": {"type": "string"},
                                "total_count": {"type": "integer"},
                                "total_amount": {"type": "number"},
                                "details": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "date": {"type": "string"},
                                            "order_count": {"type": "integer"},
                                            "total_amount": {"type": "number"},
                                            "received_amount": {"type": "number"}
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/reports/sales": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取销售报表",
                "produces": ["application/json"],
                "tags": ["报表"],
                "summary": "销售报表",
                "parameters": [
                    {"name": "start_date", "in": "query", "type": "string"},
                    {"name": "end_date", "in": "query", "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "start_date": {"type": "string"},
                                "end_date": {"type": "string"},
                                "total_count": {"type": "integer"},
                                "total_amount": {"type": "number"},
                                "details": {
                                    "type": "array",
                                    "items": {
                                        "type": "object",
                                        "properties": {
                                            "date": {"type": "string"},
                                            "order_count": {"type": "integer"},
                                            "total_amount": {"type": "number"},
                                            "delivered_amount": {"type": "number"}
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
        "/api/reports/inventory": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取库存报表",
                "produces": ["application/json"],
                "tags": ["报表"],
                "summary": "库存报表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {"$ref": "#/definitions/Inventory"}
                        }
                    }
                }
            }
        },
        "/api/contracts": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取合同列表",
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "合同列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 20}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/Contract"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建新合同",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "新建合同",
                "parameters": [
                    {"description": "合同信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Contract"}}
                ],
                "responses": {"200": {"description": "创建成功"}}
            }
        },
        "/api/contracts/{id}": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取合同详情",
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "合同详情",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/Contract"}
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新合同信息",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "编辑合同",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"},
                    {"description": "合同信息", "name": "body", "in": "body", "required": true, "schema": {"$ref": "#/definitions/Contract"}}
                ],
                "responses": {"200": {"description": "更新成功"}}
            },
            "delete": {
                "security": [{"Bearer": []}],
                "description": "删除合同",
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "删除合同",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "删除成功"}}
            }
        },
        "/api/contracts/{id}/sign": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "签署合同（甲方）",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "签署合同",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "签署成功"}}
            }
        },
        "/api/contracts/{id}/cancel": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "取消合同",
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "取消合同",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "取消成功"}}
            }
        },
        "/api/contracts/{id}/generate-token": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "生成乙方签署链接",
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "生成签署链接",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "sign_url": {"type": "string"}
                            }
                        }
                    }
                }
            }
        },
        "/api/contracts/parties": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取合同参与方列表（供应商/客户）",
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "合同参与方",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "suppliers": {"type": "array", "items": {"$ref": "#/definitions/Supplier"}},
                                "customers": {"type": "array", "items": {"$ref": "#/definitions/Customer"}}
                            }
                        }
                    }
                }
            }
        },
        "/api/public/contract/{token}": {
            "get": {
                "description": "获取公开签署链接的合同内容（无需认证）",
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "公开合同内容",
                "parameters": [
                    {"name": "token", "in": "path", "required": true, "type": "string"}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {"$ref": "#/definitions/Contract"}
                    }
                }
            }
        },
        "/api/public/contract/{token}/sign": {
            "post": {
                "description": "乙方签署合同（无需认证）",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["合同管理"],
                "summary": "乙方签署",
                "parameters": [
                    {"name": "token", "in": "path", "required": true, "type": "string"},
                    {
                        "description": "签署信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "properties": {
                                "signature": {"type": "string", "description": "签名base64"}
                            }
                        }
                    }
                ],
                "responses": {"200": {"description": "签署成功"}}
            }
        },
        "/api/system/settings": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取系统设置",
                "produces": ["application/json"],
                "tags": ["系统设置"],
                "summary": "获取系统设置",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "company_name": {"type": "string"},
                                "company_address": {"type": "string"},
                                "company_phone": {"type": "string"},
                                "currency": {"type": "string"},
                                "low_stock_threshold": {"type": "integer"},
                                "auto_backup": {"type": "string"},
                                "backup_retention": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "put": {
                "security": [{"Bearer": []}],
                "description": "更新系统设置",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["系统设置"],
                "summary": "更新系统设置",
                "parameters": [
                    {
                        "description": "设置信息",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": {"type": "string"}
                        }
                    }
                ],
                "responses": {"200": {"description": "更新成功"}}
            }
        },
        "/api/system/backups": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "获取备份记录列表",
                "produces": ["application/json"],
                "tags": ["系统设置"],
                "summary": "备份记录列表",
                "parameters": [
                    {"name": "page", "in": "query", "type": "integer", "default": 1},
                    {"name": "page_size", "in": "query", "type": "integer", "default": 50}
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"type": "array", "items": {"$ref": "#/definitions/BackupRecord"}},
                                "total": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [{"Bearer": []}],
                "description": "创建数据库备份",
                "produces": ["application/json"],
                "tags": ["系统设置"],
                "summary": "创建备份",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "data": {"$ref": "#/definitions/BackupRecord"},
                                "message": {"type": "string"}
                            }
                        }
                    }
                }
            }
        },
        "/api/system/backups/{id}/restore": {
            "post": {
                "security": [{"Bearer": []}],
                "description": "从备份恢复数据库",
                "produces": ["application/json"],
                "tags": ["系统设置"],
                "summary": "恢复备份",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "恢复成功"}}
            }
        },
        "/api/system/backups/{id}/download": {
            "get": {
                "security": [{"Bearer": []}],
                "description": "下载备份文件",
                "produces": ["application/octet-stream"],
                "tags": ["系统设置"],
                "summary": "下载备份",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "文件下载"}}
            }
        },
        "/api/system/backups/{id}": {
            "delete": {
                "security": [{"Bearer": []}],
                "description": "删除备份记录",
                "produces": ["application/json"],
                "tags": ["系统设置"],
                "summary": "删除备份",
                "parameters": [
                    {"name": "id", "in": "path", "required": true, "type": "integer"}
                ],
                "responses": {"200": {"description": "删除成功"}}
            }
        }
    },
    "definitions": {
        "Supplier": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "code": {"type": "string"},
                "name": {"type": "string"},
                "contact": {"type": "string"},
                "phone": {"type": "string"},
                "address": {"type": "string"},
                "remark": {"type": "string"},
                "status": {"type": "string"}
            }
        },
        "Customer": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "code": {"type": "string"},
                "name": {"type": "string"},
                "contact": {"type": "string"},
                "phone": {"type": "string"},
                "address": {"type": "string"},
                "remark": {"type": "string"},
                "status": {"type": "string"}
            }
        },
        "Product": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "code": {"type": "string"},
                "name": {"type": "string"},
                "category": {"type": "string"},
                "unit": {"type": "string"},
                "cost_price": {"type": "number"},
                "selling_price": {"type": "number"},
                "remark": {"type": "string"},
                "status": {"type": "string"}
            }
        },
        "PurchaseOrder": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "order_no": {"type": "string"},
                "supplier_id": {"type": "integer"},
                "supplier_name": {"type": "string"},
                "total_amount": {"type": "number"},
                "status": {"type": "string"},
                "items": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "product_id": {"type": "integer"},
                            "product_name": {"type": "string"},
                            "quantity": {"type": "integer"},
                            "unit_price": {"type": "number"},
                            "subtotal": {"type": "number"}
                        }
                    }
                },
                "created_at": {"type": "string"}
            }
        },
        "SalesOrder": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "order_no": {"type": "string"},
                "customer_id": {"type": "integer"},
                "customer_name": {"type": "string"},
                "total_amount": {"type": "number"},
                "status": {"type": "string"},
                "items": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "product_id": {"type": "integer"},
                            "product_name": {"type": "string"},
                            "quantity": {"type": "integer"},
                            "unit_price": {"type": "number"},
                            "subtotal": {"type": "number"}
                        }
                    }
                },
                "created_at": {"type": "string"}
            }
        },
        "Inventory": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "product_id": {"type": "integer"},
                "product_name": {"type": "string"},
                "quantity": {"type": "integer"},
                "cost_price": {"type": "number"}
            }
        },
        "InventoryLog": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "product_id": {"type": "integer"},
                "type": {"type": "string"},
                "quantity": {"type": "integer"},
                "balance": {"type": "integer"},
                "remark": {"type": "string"},
                "created_at": {"type": "string"}
            }
        },
        "AccountPayable": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "purchase_order_id": {"type": "integer"},
                "supplier_name": {"type": "string"},
                "amount": {"type": "number"},
                "paid_amount": {"type": "number"},
                "status": {"type": "string"},
                "created_at": {"type": "string"}
            }
        },
        "AccountReceivable": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "sales_order_id": {"type": "integer"},
                "customer_name": {"type": "string"},
                "amount": {"type": "number"},
                "paid_amount": {"type": "number"},
                "status": {"type": "string"},
                "created_at": {"type": "string"}
            }
        },
        "Expense": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "category": {"type": "string"},
                "amount": {"type": "number"},
                "remark": {"type": "string"},
                "created_at": {"type": "string"}
            }
        },
        "PaymentRecord": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "type": {"type": "string"},
                "amount": {"type": "number"},
                "payment_method": {"type": "string"},
                "remark": {"type": "string"},
                "created_at": {"type": "string"}
            }
        },
        "User": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "username": {"type": "string"},
                "real_name": {"type": "string"},
                "role": {"type": "string"}
            }
        },
        "OperationLog": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "user_id": {"type": "integer"},
                "username": {"type": "string"},
                "action": {"type": "string"},
                "target_type": {"type": "string"},
                "target_id": {"type": "integer"},
                "description": {"type": "string"},
                "ip_address": {"type": "string"},
                "created_at": {"type": "string"}
            }
        },
        "Contract": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "contract_no": {"type": "string"},
                "title": {"type": "string"},
                "type": {"type": "string"},
                "party_b_type": {"type": "string"},
                "party_b_id": {"type": "integer"},
                "party_b_name": {"type": "string"},
                "total_amount": {"type": "number"},
                "status": {"type": "string"},
                "signed_by_a": {"type": "string"},
                "signed_by_b": {"type": "string"},
                "created_at": {"type": "string"}
            }
        },
        "BackupRecord": {
            "type": "object",
            "properties": {
                "id": {"type": "integer"},
                "file_name": {"type": "string"},
                "file_size": {"type": "integer"},
                "status": {"type": "string"},
                "remark": {"type": "string"},
                "created_by": {"type": "integer"},
                "created_at": {"type": "string"}
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "JWT Bearer Token. Example: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'"
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "2.1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Buy-Demo ERP 系统 API",
	Description:      "轻量级ERP系统API文档，覆盖采购→库存→销售→财务全链路",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
