/**
 * tokenRefresher refreshes the token which is used to authenticate
 * the user, ever few minutes.
 */
app.factory('tokenRefresher', ['$interval', '$http', '$window', '$location', '$rootScope', 'toastr',
                       function($interval,   $http,   $window,   $location,   $rootScope,   toastr) {
  var refresher;
  var running = false;
  var tryOnMoreTime = 0;


  /**
   * start starts the refreshing of the token
   */
  var start = function() {
    if (!running && $window.sessionStorage.getItem('token') != undefined) {
      refresher = $interval(this.refresh, .1*60*1000);//4*60*1000); // every 4 minutes
      running = true;
    }
  };

  /**
   * stop stops the refreshing of the token. Returns true is sucessful,
   * otherwise false.
   */
  var stop = function() {
    tryOnMoreTime = 0;
    if (running) {
      running = false;
      return $interval.cancel(refresher);
    } else {
      return true;
    }
  };

  /**
   * refresh makes a one time refresh of the token
   */
  var refresh = function() {
    $http.get('/api/refreshtoken').then(
      function success(response) {
        /* The response object has these properties:

        • data – {string|Object} – The response body transformed with the transform functions.
        • status – {number} – HTTP status code of the response.
        • headers – {function([headerName])} – Header getter function.
        • config – {Object} – The configuration object that was used to generate the request.
        • statusText – {string} – HTTP status text of the response.
        */
        if (response.data.Token !== undefined) {
          /* Complete success */
          $window.sessionStorage.setItem('token', response.data.Token);

        } else {
          /* No server response */
          // try one more time

          if (tryOnMoreTime > 0) {
            toastr.error('Server error. Please log in again', {
              onTap: function() {$location.path('/login')}
            });
            $window.sessionStorage.removeItem('token');
            stop();
            $rootScope.$emit('refreshtoken-relogin'); // trigger relogin event for all subscribers

          } else {
            tryOnMoreTime++;
            refresh();
          }
        }
      },
      function error(response) {
        toastr.error('Lost connection to server. Please log in again');
        $window.sessionStorage.removeItem('token', {
              onTap: function() {$location.path('/login')}
          });
        stop();
        $rootScope.$emit('refreshtoken-relogin'); // trigger relogin event for all subscribers
      }
    )
  };

  /**
   * isRunning returns true if tokenRefresher is running, false, otherwise
   */
  var isRunning = function() {
    return running;
  };
  
  /**
   * Adds a callback to whenever the relogin event is triggered
   */
  var subscribeOnReloginEvent = function(scope, callback) {
      var handler = $rootScope.$on('refreshtoken-relogin', callback);
      scope.$on('$destroy', handler);
  };


  return {
    start: start,
    stop: stop,
    refresh: refresh,
    isRunning: isRunning,
    subscribeOnReloginEvent: subscribeOnReloginEvent
  };
}]);
