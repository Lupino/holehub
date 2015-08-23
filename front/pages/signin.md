---
layout: layout
bodyclass: signin
include_prefix: ../
---
<!-- TODO: Try to separate markup and content -->
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
      <h4>Signin</h4>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="text" id="nameOrEmail" />
        <label class="mdl-textfield__label" for="nameOrEmail">Username or Email:</label>
      </div>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="password" id="password" />
        <label class="mdl-textfield__label" for="password">Password:</label>
      </div>
    </div>
    <div class="mdl-card__actions mdl-card--border">
      <button class="mdl-button mdl-js-button mdl-button--raised mdl-button--colored" onclick="elem.signin(this);">
        Signin
      </button>
      &nbsp;
      &nbsp;
      <a href="/forget_password/index.html">
      Forget password?
      </a>
      &nbsp;
      &nbsp;
      <a href="/signup/index.html">
      I'm a new user?
      </a>
    </div>
  </div>
</section>
