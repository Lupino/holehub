---
layout: layout
bodyclass: reset_password
include_prefix: ../
---
<!-- TODO: Try to separate markup and content -->
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
      <h4>Reset Passwoed</h4>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="password" id="password" pattern=".{5}.*"/>
        <label class="mdl-textfield__label" for="password">Current password:</label>
        <span class="mdl-textfield__error">Password require more than 5 letters.</span>
      </div>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="password" id="password" pattern=".{5}.*"/>
        <label class="mdl-textfield__label" for="password">New password:</label>
        <span class="mdl-textfield__error">Password require more than 5 letters.</span>
      </div>
    </div>
    </div>
    <div class="mdl-card__actions mdl-card--border">
      <button class="mdl-button mdl-js-button mdl-button--raised mdl-button--colored">
        Reset
      </button>
    </div>
  </div>
</section>
