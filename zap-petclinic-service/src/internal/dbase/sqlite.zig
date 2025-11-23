const std = @import("std");
const c = @cImport({
    @cInclude("sqlite3.h");
});

pub const Database = struct {
    db: *c.sqlite3,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, path: []const u8) !Database {
        var db: ?*c.sqlite3 = null;
        const path_z = try allocator.dupeZ(u8, path);
        defer allocator.free(path_z);

        const result = c.sqlite3_open(path_z.ptr, &db);
        if (result != c.SQLITE_OK) {
            std.debug.print("Failed to open database: {s}\n", .{c.sqlite3_errmsg(db)});
            return error.DatabaseOpenFailed;
        }

        return Database{
            .db = db.?,
            .allocator = allocator,
        };
    }

    pub fn deinit(self: *Database) void {
        _ = c.sqlite3_close(self.db);
    }

    pub fn exec(self: *Database, sql: []const u8) !void {
        const sql_z = try self.allocator.dupeZ(u8, sql);
        defer self.allocator.free(sql_z);

        var err_msg: ?[*:0]u8 = null;
        const result = c.sqlite3_exec(self.db, sql_z.ptr, null, null, &err_msg);

        if (result != c.SQLITE_OK) {
            if (err_msg) |msg| {
                std.debug.print("SQL error: {s}\n", .{msg});
                c.sqlite3_free(msg);
            }
            return error.SqlExecFailed;
        }
    }

    pub fn createTables(self: *Database) !void {
        const create_products_table =
            \\CREATE TABLE IF NOT EXISTS products (
            \\    id INTEGER PRIMARY KEY AUTOINCREMENT,
            \\    name TEXT NOT NULL,
            \\    description TEXT,
            \\    price REAL NOT NULL,
            \\    stock INTEGER NOT NULL DEFAULT 0,
            \\    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            \\);
        ;
        try self.exec(create_products_table);
        std.debug.print("Database tables created successfully\n", .{});
    }
};
