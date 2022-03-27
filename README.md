# Instasafe
Assignment


Reqbody :- {
"amount":5000.25,
"timestamp":"{{CurrentDatetime}}"
}


Pre-Request Script:-

var moment = require('moment');
pm.globals.set("CurrentDatetime",moment().format())

