---
layout: layout
bodyclass: download
include_prefix: ../
---
<!-- TODO: Try to separate markup and content -->
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
      <h4>Download and Installation</h4>
      <p> holehub is easy to install. Download a single binary with <i>zero run-time dependencies</i> for any major platform. Unzip it and then run it from the command line.</p>
      <h5>Step 1: Download holehub</h5>
      <ul>
          <li><a href="/assets/download/darwin-amd64/holehub.zip">Mac OS X</a></li>
          <li><a href="/assets/download/darwin-386/holehub.zip">Mac OS X/i386</a></li>
          <li><a href="/assets/download/linux-amd64/holehub.zip">Linux</a></li>
          <li><a href="/assets/download/linux-arm/holehub.zip">Linux/Arm</a></li>
          <li><a href="/assets/download/linux-386/holehub.zip">Linux/i386</a></li>
          <li><a href="/assets/download/windows-arm/holehub.zip">Windows</a></li>
          <li><a href="/assets/download/windows-386/holehub.zip">Windows/i386</a></li>
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
<section class="section--center mdl-grid mdl-grid--no-spacing mdl-shadow--2dp">
  <div class="mdl-card mdl-cell mdl-cell--12-col">
    <div class="mdl-card__supporting-text">
      <h4>Build from source</h4>
      <p> holehub is easy to build from source.
        Make sure the <a href="http://golang.org"> go</a> is installed on you system.</p>
      <code>
        $ go get -v -u github.com/Lupino/holehub/holehub
      </code>
    </div>
  </div>
</section>
