const express = require("express");
const cors = require("cors");
const db = require("./database.js");
const bodyParser = require("body-parser");

const PORT = process.env.PORT || 9001;

const app = express();
app.use(cors());
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

app.get("/api/automations", (req, res) => {
  var sql = "select * from automations";
  var params = [];
  db.all(sql, params, (err, rows) => {
    if (err) {
      res.status(400).json({"error": err.message});
      return;
    }
    res.json({
      "automations": rows
    });
  });
});

app.get("/api/automations/:id", (req, res, next) => {
  var sql = "select * from automations where id = ?";
  var params = [req.params.id];
  db.get(sql, params, (err, row) => {
    if (err) {
      res.status(400).json({"error": err.message});
      return;
    }
    res.json({
      "automation": row
    });
  });
});

app.post("/api/automations", (req, res, next) => {
  var errors = [];
  if (!req.body.name) {
    errors.push("No name specified");
  }
  if (!req.body.source) {
    errors.push("No source specified");
  }
  if (!req.body.destination) {
    errors.push("No destination specified");
  }
  if (!req.body.runtime) {
    errors.push("No runtime specified");
  }
  if (!req.body.code) {
    errors.push("No code specified");
  }
  if (errors.length) {
    res.status(400).json({"error": errors.join(", ")});
    return;
  }

  var data = {
    name: req.body.name,
    source: req.body.source,
    destination: req.body.destination,
    runtime: req.body.runtime,
    code: req.body.code
  };

  var sql = "INSERT INTO automations (name, source, destination, runtime, code) VALUES (?,?,?,?,?)";
  var params = [data.name, data.source, data.destination, data.runtime, data.code];
  db.run(sql, params, function(err) {
    if (err) {
      res.status(400).json({"error": err.message});
      return;
    }
    res.json({
      "id": this.lastID
    });
  });
});

app.patch("/api/automations/:id", (req, res, next) => {
  var data = {
    name: req.body.name,
    source: req.body.source,
    destination: req.body.destination,
    runtime: req.body.runtime,
    code: req.body.code
  };

  db.run(
    `UPDATE automations SET
      name = COALESCE(?,name),
      source = COALESCE(?,source),
      destination = COALESCE(?,destination),
      runtime = COALESCE(?,runtime),
      code = COALESCE(?,code)
      WHERE ID = ?`,
    [data.name, data.source, data.destination, data.runtime, data.code, req.params.id],
    function(err) {
      if (err) {
        res.status(400).json({"error": err.message});
        return;
      }
      res.sendStatus(200);
    }
  );
});

app.delete("/api/automations/:id", (req, res, next) => {
  db.run(
    'DELETE FROM automations WHERE id = ?',
    req.params.id,
    (err) => {
      if (err) {
        res.status(400).json({"error": err.message});
        return;
      }
      res.sendStatus(200);
    }
  );
});

app.listen(PORT, () => {
  console.log(`Server listening on ${PORT}`);
});
