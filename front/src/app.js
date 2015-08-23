var request = require('superagent');
var HUB_HOST = 'http://127.0.0.1:3000';

function signin(nameOrEmail, password, callback) {
  request.post(HUB_HOST + '/api/signin/', {
    username: nameOrEmail,
    password: password
  }, function(err, res) {
    if (err) {
      return alert('Error: ' + err);
    }
    var rsp = res.body;
    if (rsp.error) {
      return alert('Error: ' + rsp.error);
    }
    window.location.href = '/signin_success/index.html';
  });
}

function signup(name, email, password, callback) {
  request.post(HUB_HOST + '/api/signup/', {
    username: name,
    email: email,
    password: password
  }, function(err, res) {
    if (err) {
      return alert('Error: ' + err);
    }
    var rsp = res.body;
    if (rsp.error) {
      return alert('Error: ' + rsp.error);
    }
    window.location.href = '/confirm_email/index.html';
  });
}

function sendResetEmail(name) {
  request.post(HUB_HOST + '/api/send/passwordToken', {
    username: name,
  }, function(err, res) {
    if (err) {
      return alert('Error: ' + err);
    }
    var rsp = res.body;
    if (rsp.error) {
      return alert('Error: ' + rsp.error);
    }
    window.location.href = '/reset_email/index.html';
  });
}


var elem = {};

elem.signin = function(e) {
  var elemNameOrEmaiil = document.getElementById('nameOrEmail');
  var elemPassword = document.getElementById('password');
  var nameOrEmail = elemNameOrEmaiil.value.trim();
  var password = elemPassword.value.trim();
  console.log(nameOrEmail, password);
  signin(nameOrEmail, password);
};

elem.signup = function(e) {
  var elemUsername = document.getElementById('username');
  var elemPassword = document.getElementById('password');
  var elemEmail = document.getElementById('email');
  var userName = elemUsername.value.trim();
  var password = elemPassword.value.trim();
  var email = elemEmail.value.trim();
  signup(userName, email, password);
};

elem.sendResetEmail = function(e) {
  var elemUsername = document.getElementById('username');
  var username = elemUsername.value.trim();
  sendResetEmail(username);
};

window['elem'] = elem;
