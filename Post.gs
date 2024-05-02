function onSubmit(e) {

  var form = FormApp.getActiveForm();
  var formResponses = form.getResponses()
  var data = new Object();

  for (var i = 0; i < formResponses.length; i++) {
    var formResponse = formResponses[i];
    var itemResponses = formResponse.getItemResponses();
  }

  data.name = itemResponses[0].getResponse();
  data.email = itemResponses[1].getResponse();
  data.message = itemResponses[2].getResponse();
  data.file = 'https://drive.google.com/open?id=' + String(itemResponses[3].getResponse()[0]);

  var options = {
    'method' : 'post',
    'contentType': 'application/json',
    // Convert the JavaScript object to a JSON string.
    'payload' : JSON.stringify(data)
  };

  UrlFetchApp.fetch('https://formapi.fly.dev/sheet', options);
}
