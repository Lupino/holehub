---
layout: layout
bodyclass: overview
include_prefix: ./
---
<!-- TODO: Try to separate markup and content -->
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <header class="section__play-btn mdl-cell mdl-cell--3-col-desktop mdl-cell--2-col-tablet mdl-cell--4-col-phone mdl-color--teal-100 mdl-color-text--white">
    <i class="material-icons">play_circle_filled</i>
  </header>
  <div class="mdl-card mdl-cell mdl-cell--9-col-desktop mdl-cell--6-col-tablet mdl-cell--4-col-phone">
    <div class="mdl-card__supporting-text">
      <h4>Secure tunnels to localhost</h4>
      <p>"I want to expose a local server behind a NAT or firewall to the internet."</p>
      <!-- Colored raised button -->
      <a href="/signin/index.html" class="mdl-button mdl-js-button mdl-button--raised mdl-button--colored">
        Signin
      </a> &nbsp; Or &nbsp;
      <a href="/signup/index.html" class="mdl-button mdl-js-button mdl-button--raised mdl-button--accent">
        Signup
      </a>
    </div>
  </div>
</section>
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text mdl-grid mdl-grid--no-spacing">
      <h4 class="mdl-cell mdl-cell--12-col">Features</h4>
      <div class="section__circle-container mdl-cell mdl-cell--2-col mdl-cell--1-col-phone">
        <div class="section__circle-container__circle mdl-color--primary"></div>
      </div>
      <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
        <h5><b>Demo without deploying</b></h5>
        <p>Don't constantly redeploy your in-progress work to get feedback from clients.
        holehub creates a secure public URL (http://holehub.com:port) to a local webserver on you machine.
        Iterate quickly with immediate feedback without interrupting flow.</p>
      </div>
      <div class="section__circle-container mdl-cell mdl-cell--2-col mdl-cell--1-col-phone">
        <div class="section__circle-container__circle mdl-color--primary"></div>
      </div>
      <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
        <h5><b>Simplify mobile device testing</b></h5>
        <p>Test mobile apps against a development backend running on your machine.
        Point holehub at your local dev server and then configure your app to use the holehub URL.
        It won't change, event when you change networks.</p>
      </div>
      <div class="section__circle-container mdl-cell mdl-cell--2-col mdl-cell--1-col-phone">
        <div class="section__circle-container__circle mdl-color--primary"></div>
      </div>
      <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
        <h5><b>Build webhook consumers with ease</b></h5>
        <p>Building webhook consumers can be a pain: it requires a public address and a lot of set up to trigger hooks.
        Save yourself time and frustration with holehub. Inspect the HTTP traffic flowing over your tunnel.
        Then replay webhooks requests with one click to iterate quickly while staying in context.</p>
      </div>
      <div class="section__circle-container mdl-cell mdl-cell--2-col mdl-cell--1-col-phone">
        <div class="section__circle-container__circle mdl-color--primary"></div>
      </div>
      <div class="section__text mdl-cell mdl-cell--10-col-desktop mdl-cell--6-col-tablet mdl-cell--3-col-phone">
        <h5><b>Run personal cloud services from your own private network</b></h5>
        <p>Own your data. Host personal cloud services on your own private network. Run webmail, file syncing, and more securely on your hardware with full end-to-end encryotion.</p>
      </div>
    </div>
  </div>
</section>
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
    <p>Need to run <a> holehub in production</a>? </p>
    <p>Use <a>holehub link</a> to manage your IoT devices or as a lightweight alternative to VPN for targeted access into customer networks.</p>
    </div>
  </div>
</section>
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
      <h4>Download and Installation</h4>
      <p> holehub is easy to install. Download a single binary with <i>zero run-time dependencies</i> for any major platform. Unzip it and then run it from the command line.</p>
      <h5>Step 1: Download holehub</h5>
      <ul>
          <li><a href="/assets/download/darwin-amd64/holehub.zip">Mac OS X</a></li>
          <li><a href="/assets/download/linux-amd64/holehub.zip">Linux</a></li>
          <li><a href="/assets/download/linux-arm/holehub.zip">Linux/Arm</a></li>
          <li><a href="/assets/download/windows-amd64/holehub.zip">Windows</a></li>
          <li><a href="/download/index.html">More platforms</a></li>
      </ul>
      <h5>Setup 2: Unzip it</h5>
      <p>On linux or OSX you can unzip holehub from a terminal with the following command. On windows, just double click holehub.zip </p>
      <code>
      $ unzip path/to/holehub.zip
      </code>

      <h5>Setup 3: Run it!</h5>
      <p> Read the <a href="/docs/index.html"> Usage Guide</a> for documentation on how to use holehub. Try it out running it from the command line:</p>
      <code>
        $ ./holehub --help
      </code>
    </div>
  </div>
</section>
