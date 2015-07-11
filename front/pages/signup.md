---
layout: layout
bodyclass: signup
include_prefix: ../
---
<!-- TODO: Try to separate markup and content -->
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
      <h4>Signup</h4>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="text" id="username" pattern="[a-zA-Z][0-9a-zA-Z]{4}[0-9a-zA-Z]*" />
        <label class="mdl-textfield__label" for="username">Username:</label>
        <span class="mdl-textfield__error">Letters or number only</span>
      </div>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="text" id="email" pattern="(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})"/>
        <label class="mdl-textfield__label" for="email">Email:</label>
        <span class="mdl-textfield__error">InvalId email address</span>
      </div>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="password" id="password" pattern=".{5}.*"/>
        <label class="mdl-textfield__label" for="password">Password:</label>
        <span class="mdl-textfield__error">Password require more than 5 letters.</span>
      </div>
    </div>
    <div class="mdl-card__actions mdl-card--border">
      <button class="mdl-button mdl-js-button mdl-button--raised mdl-button--colored">
        Signup
      </button>
      &nbsp;
      &nbsp;
      <a href="/signin/index.html">
      I'm an old user?
      </a>
    </div>
  </div>
</section>
