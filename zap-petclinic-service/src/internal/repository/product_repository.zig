const std = @import("std");
const Product = @import("../models/product.zig").Product;
const CreateProductRequest = @import("../models/product.zig").CreateProductRequest;
const UpdateProductRequest = @import("../models/product.zig").UpdateProductRequest;
const Database = @import("../database/sqlite.zig").Database;
const c = @cImport({
    @cInclude("sqlite3.h");
});

pub const ProductRepository = struct {
    db: *Database,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, db: *Database) ProductRepository {
        return ProductRepository{
            .db = db,
            .allocator = allocator,
        };
    }

    pub fn create(self: *ProductRepository, req: CreateProductRequest) !i64 {
        const sql = "INSERT INTO products (name, description, price, stock) VALUES (?, ?, ?, ?)";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        const name_z = try self.allocator.dupeZ(u8, req.name);
        defer self.allocator.free(name_z);
        const desc_z = try self.allocator.dupeZ(u8, req.description);
        defer self.allocator.free(desc_z);

        _ = c.sqlite3_bind_text(stmt, 1, name_z.ptr, -1, c.SQLITE_TRANSIENT);
        _ = c.sqlite3_bind_text(stmt, 2, desc_z.ptr, -1, c.SQLITE_TRANSIENT);
        _ = c.sqlite3_bind_double(stmt, 3, req.price);
        _ = c.sqlite3_bind_int(stmt, 4, req.stock);

        result = c.sqlite3_step(stmt);
        if (result != c.SQLITE_DONE) {
            return error.InsertFailed;
        }

        return c.sqlite3_last_insert_rowid(self.db.db);
    }

    pub fn findAll(self: *ProductRepository) !std.ArrayList(Product) {
        const sql = "SELECT id, name, description, price, stock, created_at FROM products";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        var products = std.ArrayList(Product).init(self.allocator);

        while (c.sqlite3_step(stmt) == c.SQLITE_ROW) {
            const id = c.sqlite3_column_int64(stmt, 0);
            const name_ptr = c.sqlite3_column_text(stmt, 1);
            const desc_ptr = c.sqlite3_column_text(stmt, 2);
            const price = c.sqlite3_column_double(stmt, 3);
            const stock = c.sqlite3_column_int(stmt, 4);
            const created_ptr = c.sqlite3_column_text(stmt, 5);

            const name = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(name_ptr))));
            const desc = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(desc_ptr))));
            const created = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(created_ptr))));

            try products.append(Product{
                .id = id,
                .name = name,
                .description = desc,
                .price = price,
                .stock = stock,
                .created_at = created,
            });
        }

        return products;
    }

    pub fn findById(self: *ProductRepository, id: i64) !?Product {
        const sql = "SELECT id, name, description, price, stock, created_at FROM products WHERE id = ?";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        _ = c.sqlite3_bind_int64(stmt, 1, id);

        if (c.sqlite3_step(stmt) == c.SQLITE_ROW) {
            const name_ptr = c.sqlite3_column_text(stmt, 1);
            const desc_ptr = c.sqlite3_column_text(stmt, 2);
            const price = c.sqlite3_column_double(stmt, 3);
            const stock = c.sqlite3_column_int(stmt, 4);
            const created_ptr = c.sqlite3_column_text(stmt, 5);

            const name = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(name_ptr))));
            const desc = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(desc_ptr))));
            const created = try self.allocator.dupe(u8, std.mem.span(@as([*:0]const u8, @ptrCast(created_ptr))));

            return Product{
                .id = id,
                .name = name,
                .description = desc,
                .price = price,
                .stock = stock,
                .created_at = created,
            };
        }

        return null;
    }

    pub fn deleteById(self: *ProductRepository, id: i64) !bool {
        const sql = "DELETE FROM products WHERE id = ?";
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var stmt: ?*c.sqlite3_stmt = null;
        var result = c.sqlite3_prepare_v2(self.db.db, sql_z.ptr, -1, &stmt, null);

        if (result != c.SQLITE_OK) {
            return error.PrepareFailed;
        }
        defer _ = c.sqlite3_finalize(stmt);

        _ = c.sqlite3_bind_int64(stmt, 1, id);
        result = c.sqlite3_step(stmt);

        if (result != c.SQLITE_DONE) {
            return error.DeleteFailed;
        }

        return c.sqlite3_changes(self.db.db) > 0;
    }
};
