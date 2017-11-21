var levelgraph = require("levelgraph");
var level = require("level-mem");

var lvl = level("test");
var db = levelgraph(lvl);

var stream = db.putStream();

db.put(
  [
    {
      subject: "matteo",
      predicate: "friend",
      object: "daniele"
    },
    {
      subject: "daniele",
      predicate: "friend",
      object: "matteo"
    },
    {
      subject: "daniele",
      predicate: "friend",
      object: "marco"
    },
    {
      subject: "lucio",
      predicate: "friend",
      object: "matteo"
    },
    {
      subject: "lucio",
      predicate: "friend",
      object: "marco"
    },
    {
      subject: "marco",
      predicate: "friend",
      object: "davide"
    }
  ],
  function() {
    var stream = db.searchStream([
      {
        subject: "matteo",
        predicate: "friend",
        object: db.v("x")
      },
      {
        subject: db.v("x"),
        predicate: "friend",
        object: db.v("y")
      },
      {
        subject: db.v("y"),
        predicate: "friend",
        object: "davide"
      }
    ]);

    stream.on("data", function(data) {
      // this will print "{ x: 'daniele', y: 'marco' }"
      console.log(data);
    });
  }
);
