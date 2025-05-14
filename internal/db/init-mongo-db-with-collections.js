db.createCollection("users", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name", "email", "password"],
            properties: {
                name: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                email: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                password: {
                    bsonType: "string",
                    description: "must be a string and is required"
                }
            }
        }
    }
});

db.createCollection("tasks", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["title", "userId"],
            properties: {
                title: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                description: {
                    bsonType: "string",
                    description: "must be a string"
                },
                completed: {
                    bsonType: "bool",
                    description: "must be a boolean"
                },
                userId: {
                    bsonType: "objectId",
                    description: "must be an objectId and is required"
                }
            }
        }
    }
});

