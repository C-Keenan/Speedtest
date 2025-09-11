# Speedtest
This repository is for housing the code and scripts associated with my automated speedtest logging containers (currently un-published) and should you wish to download, compile and run them you may wish to obtain your own copy of the 'Speedtest CLI' binary, you can find said binary at the link below.

I plan to reach out to Ookla or the subsidiary Speedtest.net(?) to ascertain whether or not I can and should publish the containers but that will likely take some time depending on their potential demands and/or required changes. Depending on their answer I will be providing a non-build version of the docker-compose.yml file.

If you have any questions please reach out to dev+speedtest-questions@lollypopstealer.com and I will do my best to answer them in a timely manner. Any attempt to reach out directly to the main inbox will be ignored and deleted, any spam found in the speedtest-questions folder will be reported and your email address added to my ignore list. Please be respectful to me and my email inbox.

I am self-taught and still learning to code and I do utilize AI a bit to troubleshoot or provide copy/paste output from time to time as my wrists seize up at the best of times from typing too much. All code is tested manually the same way I would expect it to be run (in docker, in this case) and has NOT been tested in any other scenarios. Feel free to test your preferred scenario and report back your findings. I will keep a list of tested run scenarios in this readme below the download link for the speedtest binary. You can report your tests to dev+speedtest-runtime-tests@lollypopstealer.com with a youtube video link, link shorteners will be ignored and all rules from above also apply here, confirming the test run and a detailed list of steps to repeat the process as well as what you'd like to be reffered to as so I can give you credit. I will verify your findings myself and add your name as a tester and as contributer if your code (should you submit any (also put in the email)) leads towards stability or other improvements.

## Speedtest binary download link
[Download the Ookla Speedtest CLI Binary for your system](https://www.speedtest.net/apps/cli)

## Updates
*Nothing to see hear yet*

## Future Plans
- Add a button on the webpage to manually run a speedtest and when finished reload the page
- Show the average daily speed and ping times in a small table at the top (helps evaluate ISP downtime or hiccups at a glance)

## Tested Run Scenarios
- Docker: Tested, Running
- WSL - Ubuntu - 24.04: Tested, Running
- Bare Metal: Un-tested
- Kubernetes: Un-tested
