db = db.getSiblingDB('oe');
db.createUser({
  user: 'oe',
  pwd: 'oe',
  roles: [{
    role: 'readWrite',
    db: 'oe'
  }]
})

db.createCollection('partner', { capped: false });

db.partner.insert([
  { "name": "Alvabot 1", "rating": 3, "orders": 5 },
  { "name": "Alvabot 2", "rating": 4, "orders": 20 },
  { "name": "Alvabot 3", "rating": 5, "orders": 15 },
  { "name": "Alvabot 4", "rating": 5, "orders": 5 },
  { "name": "Alvabot 5", "rating": 2, "orders": 5 },
  { "name": "Alvabot 6", "rating": 5, "orders": 10 },
]);