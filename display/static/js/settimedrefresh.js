function setTimedRefresh() {
  const now = new Date();
  const nextRefresh = new Date();
            
  nextRefresh.setHours(now.getHours() + 1);
  nextRefresh.setMinutes(10);
  nextRefresh.setSeconds(0);
  nextRefresh.setMilliseconds(0);

  if (nextRefresh.getTime() <= now.getTime()) {
    nextRefresh.setHours(nextRefresh.getHours() + 1);
  }

  const timeout = nextRefresh.getTime() - now.getTime();
            
  console.log(`Page will refresh in ${timeout / 1000} seconds at ${nextRefresh.toLocaleTimeString()}`);
            
  setTimeout(() => {
    window.location.reload(true);
  }, timeout);
}

setTimedRefresh();