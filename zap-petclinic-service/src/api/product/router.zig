const std = @import("std");
const zap = @import("zap");
const Product = @import("../models/product.zig").Product;
const CreateProductRequest = @import("../models/product.zig").CreateProductRequest;
const ProductRepository = @import("../repositories/product_repository.zig").ProductRepository;

pub const ProductHandler = struct {
    repo: *ProductRepository,
    allocator: std.mem.Allocator,

    pub fn init(allocator: std.mem.Allocator, repo: *ProductRepository) ProductHandler {
        return ProductHandler{
            .repo = repo,
            .allocator = allocator,
        };
    }

    pub fn getAllProducts(self: *ProductHandler, r: zap.Request) void {
        var products = self.repo.findAll() catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to fetch products\"}") catch return;
            return;
        };
        defer {
            for (products.items) |*p| {
                p.deinit(self.allocator);
            }
            products.deinit();
        }

        var response = std.ArrayList(u8).init(self.allocator);
        defer response.deinit();

        response.appendSlice("{\"products\":[") catch return;
        for (products.items, 0..) |product, i| {
            if (i > 0) response.appendSlice(",") catch return;
            std.fmt.format(response.writer(), "{{\"id\":{d},\"name\":\"{s}\",\"description\":\"{s}\",\"price\":{d:.2},\"stock\":{d},\"created_at\":\"{s}\"}}", .{
                product.id,
                product.name,
                product.description,
                product.price,
                product.stock,
                product.created_at,
            }) catch return;
        }
        response.appendSlice("]}") catch return;

        r.setStatus(.ok);
        r.setContentType(.JSON) catch return;
        r.sendBody(response.items) catch return;
    }

    pub fn getProductById(self: *ProductHandler, r: zap.Request) void {
        const path = r.path orelse {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid path\"}") catch return;
            return;
        };

        var it = std.mem.splitSequence(u8, path, "/");
        var id_str: ?[]const u8 = null;
        var count: usize = 0;
        while (it.next()) |segment| {
            count += 1;
            if (count == 3) id_str = segment;
        }

        const id = std.fmt.parseInt(i64, id_str orelse "", 10) catch {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid product ID\"}") catch return;
            return;
        };

        var product = self.repo.findById(id) catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to fetch product\"}") catch return;
            return;
        };

        if (product) |*p| {
            defer p.deinit(self.allocator);

            var response = std.ArrayList(u8).init(self.allocator);
            defer response.deinit();

            std.fmt.format(response.writer(), "{{\"id\":{d},\"name\":\"{s}\",\"description\":\"{s}\",\"price\":{d:.2},\"stock\":{d},\"created_at\":\"{s}\"}}", .{
                p.id,
                p.name,
                p.description,
                p.price,
                p.stock,
                p.created_at,
            }) catch return;

            r.setStatus(.ok);
            r.setContentType(.JSON) catch return;
            r.sendBody(response.items) catch return;
        } else {
            r.setStatus(.not_found);
            r.sendBody("{\"error\":\"Product not found\"}") catch return;
        }
    }

    pub fn createProduct(self: *ProductHandler, r: zap.Request) void {
        const body = r.body orelse {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Request body required\"}") catch return;
            return;
        };

        const parsed = std.json.parseFromSlice(CreateProductRequest, self.allocator, body, .{}) catch {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid JSON\"}") catch return;
            return;
        };
        defer parsed.deinit();

        const id = self.repo.create(parsed.value) catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to create product\"}") catch return;
            return;
        };

        var response = std.ArrayList(u8).init(self.allocator);
        defer response.deinit();

        std.fmt.format(response.writer(), "{{\"id\":{d},\"message\":\"Product created successfully\"}}", .{id}) catch return;

        r.setStatus(.created);
        r.setContentType(.JSON) catch return;
        r.sendBody(response.items) catch return;
    }

    pub fn deleteProduct(self: *ProductHandler, r: zap.Request) void {
        const path = r.path orelse {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid path\"}") catch return;
            return;
        };

        var it = std.mem.splitSequence(u8, path, "/");
        var id_str: ?[]const u8 = null;
        var count: usize = 0;
        while (it.next()) |segment| {
            count += 1;
            if (count == 3) id_str = segment;
        }

        const id = std.fmt.parseInt(i64, id_str orelse "", 10) catch {
            r.setStatus(.bad_request);
            r.sendBody("{\"error\":\"Invalid product ID\"}") catch return;
            return;
        };

        const deleted = self.repo.deleteById(id) catch {
            r.setStatus(.internal_server_error);
            r.sendBody("{\"error\":\"Failed to delete product\"}") catch return;
            return;
        };

        if (deleted) {
            r.setStatus(.ok);
            r.sendBody("{\"message\":\"Product deleted successfully\"}") catch return;
        } else {
            r.setStatus(.not_found);
            r.sendBody("{\"error\":\"Product not found\"}") catch return;
        }
    }
};
