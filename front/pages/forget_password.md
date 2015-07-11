---
layout: layout
bodyclass: forget_password
include_prefix: ../
---
<!-- TODO: Try to separate markup and content -->
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
      <h4>Forget Password</h4>
      <div class="mdl-textfield mdl-js-textfield mdl-textfield--floating-label">
        <input class="mdl-textfield__input" type="text" id="username" />
        <label class="mdl-textfield__label" for="username">Username:</label>
      </div>
      <button class="mdl-button mdl-js-button mdl-button--primary">
        Send
      </button>
      &nbsp;
      &nbsp;
      <a href="/signin/index.html">
      I remember my password?
      </a>
  </div>
</section>
