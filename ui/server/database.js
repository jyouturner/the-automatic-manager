var sqlite3 = require("sqlite3").verbose()

const DBSOURCE = "./server/db.sqlite";

let db = new sqlite3.Database(DBSOURCE, (err) => {
  if (err) {
    console.error(err.message);
    throw err;
  } else {
    console.log("Connected to the SQLite database.");
    db.run(`CREATE TABLE automations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name text NOT NULL,
        source text NOT NULL,
        destination text NOT NULL,
        runtime text NOT NULL,
        code blob DEFAULT ""
      )`,
    (err) => {
      if (err) {
        // Table already created
      } else {
        // Table just created, inserting some dummy data
        var insert = `INSERT INTO automations (name, source, destination, runtime) VALUES (?,?,?,?)`;
        db.run(insert, ["Automation 1", "Jira", "Notion", "Node"]);
        db.run(insert, ["Automation 2", "Confluence", "Notion", "Python"]);
        db.run(insert, ["Automation 3", "Notion", "Slack", "Node"]);
      }
    });
  }
});

module.exports = db;
