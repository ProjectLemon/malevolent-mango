/**
 * tokenRefresher refreshes the token which is used to authenticate
 * the user, ever few minutes.
 */
app.factory('tokenRefresher', ['$interval', '$http', '$window', 
                       function($interval,   $http,   $window) {
  var refresher;
  var running = false;
  
  return {    
    /**
     * start starts the refreshing of the token
     */
    start: function() {
      if (!running && $window.sessionStorage.token !== undefined) {
        refresher = $interval(this.refresh, 4*60*1000); // every 4 minutes
        running = true;
      }
    },
  
    /**
     * stop stops the refreshing of the token. Returns true is sucessful,
     * otherwise false.
     */  
    stop: function() {
      if (running) {
        running = false;
        return $interval.cancel(refresher);
      } else {
        return true;
      }
    },
  
    /**
     * refresh makes a one time refresh of the token
     */
    refresh: function() {
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
            $window.sessionStorage.token = response.data.Token;
            
          } else {
            /* No server response */
          }
        },
        function error(response) {
          console.log('Error: could not refresh token ('+response.status+')');
        }
      )
    },
    
    /**
     * isRunning returns true is tokenRefresher, false, otherwise
     */
    isRunning: function() {
      return running;
    }
  };
}]);