const std = @import("std");
const zap = @import("zap");
const sqlite = @import("sqlite");

const Product = struct {
    id: i64,
    name: []const u8,
    description: []const u8,
    price: f64,
    stock: i32,
};

const ProductInput = struct {
    name: []const u8,
    description: []const u8,
    price: f64,
    stock: i32,
};

var db: sqlite.Db = undefined;
var gpa = std.heap.GeneralPurposeAllocator(.{}){};

fn initDatabase() !void {
    db = try sqlite.Db.init(.{
        .mode = sqlite.Db.Mode{ .File = "products.db" },
        .open_flags = .{
            .write = true,
            .create = true,
        },
    });

    const create_table =
        \\CREATE TABLE IF NOT EXISTS products (
        \\  id INTEGER PRIMARY KEY AUTOINCREMENT,
        \\  name TEXT NOT NULL,
        \\  description TEXT,
        \\  price REAL NOT NULL,
        \\  stock INTEGER NOT NULL DEFAULT 0
        \\)
    ;

    try db.exec(create_table, .{}, .{});
}

fn getAllProducts(r: zap.Request) void {
    const allocator = gpa.allocator();

    var stmt = db.prepareDynamic(
        "SELECT id, name, description, price, stock FROM products"
    ) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Database error") catch return;
        return;
    };
    defer stmt.deinit();

    var products = std.ArrayList(Product).init(allocator);
    defer {
        for (products.items) |p| {
            allocator.free(p.name);
            allocator.free(p.description);
        }
        products.deinit();
    }

    var iter = stmt.iterator(Product, .{}) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Query error") catch return;
        return;
    };

    while (iter.next(.{}) catch null) |row| {
        const product = Product{
            .id = row.id,
            .name = allocator.dupe(u8, row.name) catch continue,
            .description = allocator.dupe(u8, row.description) catch continue,
            .price = row.price,
            .stock = row.stock,
        };
        products.append(product) catch continue;
    }

    var json = std.ArrayList(u8).init(allocator);
    defer json.deinit();

    std.json.stringify(products.items, .{}, json.writer()) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("JSON error") catch return;
        return;
    };

    r.setHeader("Content-Type", "application/json") catch return;
    r.sendBody(json.items) catch return;
}

fn getProduct(r: zap.Request) void {
    const allocator = gpa.allocator();

    const id_str = r.getParamStr("id") orelse {
        r.setStatus(.bad_request);
        r.sendBody("Missing ID parameter") catch return;
        return;
    };

    const id = std.fmt.parseInt(i64, id_str.str, 10) catch {
        r.setStatus(.bad_request);
        r.sendBody("Invalid ID") catch return;
        return;
    };

    var stmt = db.prepareDynamic(
        "SELECT id, name, description, price, stock FROM products WHERE id = ?"
    ) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Database error") catch return;
        return;
    };
    defer stmt.deinit();

    const row = stmt.oneAlloc(Product, allocator, .{}, .{ .id = id }) catch |err| {
        if (err == error.NoRows) {
            r.setStatus(.not_found);
            r.sendBody("Product not found") catch return;
            return;
        }
        r.setStatus(.internal_server_error);
        r.sendBody("Query error") catch return;
        return;
    };

    if (row) |product| {
        defer {
            allocator.free(product.name);
            allocator.free(product.description);
        }

        var json = std.ArrayList(u8).init(allocator);
        defer json.deinit();

        std.json.stringify(product, .{}, json.writer()) catch {
            r.setStatus(.internal_server_error);
            r.sendBody("JSON error") catch return;
            return;
        };

        r.setHeader("Content-Type", "application/json") catch return;
        r.sendBody(json.items) catch return;
    }
}

fn createProduct(r: zap.Request) void {
    const allocator = gpa.allocator();

    const body = r.body orelse {
        r.setStatus(.bad_request);
        r.sendBody("Missing request body") catch return;
        return;
    };

    const parsed = std.json.parseFromSlice(
        ProductInput,
        allocator,
        body,
        .{}
    ) catch {
        r.setStatus(.bad_request);
        r.sendBody("Invalid JSON") catch return;
        return;
    };
    defer parsed.deinit();

    const input = parsed.value;

    var stmt = db.prepareDynamic(
        "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)"
    ) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Database error") catch return;
        return;
    };
    defer stmt.deinit();

    stmt.exec(.{}, .{
        .name = input.name,
        .description = input.description,
        .price = input.price,
        .stock = input.stock,
    }) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Insert failed") catch return;
        return;
    };

    const last_id = db.getLastInsertRowId();

    var response = std.ArrayList(u8).init(allocator);
    defer response.deinit();

    std.fmt.format(response.writer(), "{{\"id\":{d},\"message\":\"Product created\"}}", .{last_id}) catch {
        r.setStatus(.internal_server_error);
        return;
    };

    r.setStatus(.created);
    r.setHeader("Content-Type", "application/json") catch return;
    r.sendBody(response.items) catch return;
}

fn updateProduct(r: zap.Request) void {
    const allocator = gpa.allocator();

    const id_str = r.getParamStr("id") orelse {
        r.setStatus(.bad_request);
        r.sendBody("Missing ID parameter") catch return;
        return;
    };

    const id = std.fmt.parseInt(i64, id_str.str, 10) catch {
        r.setStatus(.bad_request);
        r.sendBody("Invalid ID") catch return;
        return;
    };

    const body = r.body orelse {
        r.setStatus(.bad_request);
        r.sendBody("Missing request body") catch return;
        return;
    };

    const parsed = std.json.parseFromSlice(
        ProductInput,
        allocator,
        body,
        .{}
    ) catch {
        r.setStatus(.bad_request);
        r.sendBody("Invalid JSON") catch return;
        return;
    };
    defer parsed.deinit();

    const input = parsed.value;

    var stmt = db.prepareDynamic(
        "UPDATE products SET name = ?, description = ?, price = ?, stock = ? WHERE id = ?"
    ) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Database error") catch return;
        return;
    };
    defer stmt.deinit();

    stmt.exec(.{}, .{
        .name = input.name,
        .description = input.description,
        .price = input.price,
        .stock = input.stock,
        .id = id,
    }) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Update failed") catch return;
        return;
    };

    if (db.getLastChangesCount() == 0) {
        r.setStatus(.not_found);
        r.sendBody("Product not found") catch return;
        return;
    };

    r.setHeader("Content-Type", "application/json") catch return;
    r.sendBody("{\"message\":\"Product updated\"}") catch return;
}

fn deleteProduct(r: zap.Request) void {
    const id_str = r.getParamStr("id") orelse {
        r.setStatus(.bad_request);
        r.sendBody("Missing ID parameter") catch return;
        return;
    };

    const id = std.fmt.parseInt(i64, id_str.str, 10) catch {
        r.setStatus(.bad_request);
        r.sendBody("Invalid ID") catch return;
        return;
    };

    var stmt = db.prepareDynamic(
        "DELETE FROM products WHERE id = ?"
    ) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Database error") catch return;
        return;
    };
    defer stmt.deinit();

    stmt.exec(.{}, .{ .id = id }) catch {
        r.setStatus(.internal_server_error);
        r.sendBody("Delete failed") catch return;
        return;
    };

    if (db.getLastChangesCount() == 0) {
        r.setStatus(.not_found);
        r.sendBody("Product not found") catch return;
        return;
    };

    r.setHeader("Content-Type", "application/json") catch return;
    r.sendBody("{\"message\":\"Product deleted\"}") catch return;
}

pub fn main() !void {
    try initDatabase();
    defer db.deinit();

    var listener = zap.HttpListener.init(.{
        .port = 3000,
        .on_request = null,
        .log = true,
    });

    var router = zap.Router.init(gpa.allocator(), .{});
    defer router.deinit();

    try router.handle_func("/api/products", getAllProducts, .GET);
    try router.handle_func("/api/products/:id", getProduct, .GET);
    try router.handle_func("/api/products", createProduct, .POST);
    try router.handle_func("/api/products/:id", updateProduct, .PUT);
    try router.handle_func("/api/products/:id", deleteProduct, .DELETE);

    listener.on_request = router.on_request_handler();

    try listener.listen();

    std.debug.print("Server running on http://localhost:3000\n", .{});
    std.debug.print("Endpoints:\n", .{});
    std.debug.print("  GET    /api/products     - Get all products\n", .{});
    std.debug.print("  GET    /api/products/:id - Get product by ID\n", .{});
    std.debug.print("  POST   /api/products     - Create product\n", .{});
    std.debug.print("  PUT    /api/products/:id - Update product\n", .{});
    std.debug.print("  DELETE /api/products/:id - Delete product\n", .{});

    zap.start(.{
        .threads = 2,
        .workers = 2,
    });
}