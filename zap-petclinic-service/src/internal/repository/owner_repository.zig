const OwnerRepositorier = struct {

};









const OwnerRepository = struct {
    const std = @import("std");
    const Owner = @import("../model/owner.zig").Owner;
    const Database = @import("../database/database.zig").Database;

    db: *Database,

    pub fn init(db: *Database) OwnerRepository {
        return OwnerRepository{
            .db = db,
        };
    }

    pub fn findById(self: *OwnerRepository, id: i32) ?Owner {
        const query = "SELECT id, first_name, last_name, address, city, telephone FROM owners WHERE id = ?";
        var stmt = self.db.prepare(query) catch return null;
        defer stmt.finalize();

        stmt.bindInt(1, id) catch return null;

        if (stmt.step() == .ROW) {
            return Owner{
                .id = stmt.columnInt(0),
                .first_name = stmt.columnText(1),
                .last_name = stmt.columnText(2),
                .address = stmt.columnText(3),
                .city = stmt.columnText(4),
                .telephone = stmt.columnText(5),
            };
        } else {
            return null;
        }
    }

    pub fn save(self: *OwnerRepository, owner: Owner) !void {
        const query = "INSERT INTO owners (first_name, last_name, address, city, telephone) VALUES (?, ?, ?, ?, ?)";
        var stmt = self.db.prepare(query) catch return error.DatabaseError;
        defer stmt.finalize();

        stmt.bindText(1, owner.first_name) catch return error.DatabaseError;
        stmt.bindText(2, owner.last_name) catch return error.DatabaseError;
        stmt.bindText(3, owner.address) catch return error.DatabaseError;
        stmt.bindText(4, owner.city) catch return error.DatabaseError;
        stmt.bindText(5, owner.telephone) catch return error.DatabaseError;

        if (stmt.step() != .DONE) {
            return error.DatabaseError;
        }
    }
};