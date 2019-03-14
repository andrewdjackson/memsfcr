
var spawn = require('child_process').spawn;
var readmems = spawn('readmems',  ['/dev/cu.usbserial-FT94CQQS', 'interactive']);

readmems.stdout.on('data', (data) => {
  console.log(`stdout: ${data}`);
});

readmems.stderr.on('data', (data) => {
  console.log(`stderr: ${data}`);
});

readmems.on('close', (code) => {
  console.log(`child process exited with code ${code}`);
});


// readmems.stdin.setEncoding('utf-8');
// readmems.stdin.write("console.log('Hello from PhantomJS')\n");
// readmems.stdin.end();