describe('Token Refresher', function() {
  
  var tokenRefresher;
  var $httpBackend;
  var $interval;
  var $window;
  
  var token = 1;
  var refreshTokenURL = '/api/refreshtoken';
  
  beforeEach(module('MaliciousMango'));
  beforeEach(inject(function($injector, _$interval_, _$window_, _tokenRefresher_) {
    // Mock up fake server responses
    $httpBackend = $injector.get('$httpBackend');
    $httpBackend.when('GET', refreshTokenURL)
      .respond(function(method, url, data, headers) {
        token++; 
        return [200, {'Token': token}, {}]; // status, body, header
      });
    
    // Load angular methods
    $interval = _$interval_;
    $window = _$window_;
    $window.sessionStorage.setItem('token', token);
    // Load factory
    tokenRefresher = _tokenRefresher_;
    tokenRefresher.stop();
  }));
  
  /* Make sure no other request was called */
  afterEach(function() {
    $httpBackend.verifyNoOutstandingExpectation();
    $httpBackend.verifyNoOutstandingRequest();
    tokenRefresher.stop();
  });

  
  it('should be instantiated', function() {
    expect(tokenRefresher).toBeDefined();
  });  
  
  it('updates token after 4 minutes', function() {
    $httpBackend.expectGET(refreshTokenURL);
    tokenRefresher.start();
    $interval.flush(4*60*1000+10);
    $httpBackend.flush();
    expect($window.sessionStorage.token).toBe(token.toString());
    tokenRefresher.stop();
  });
  
  it('updates token every 4 minutes (12 times)', function() {
    tokenRefresher.start();
    for (i = 0; i < 12; i++) {
      $httpBackend.expectGET(refreshTokenURL);
      $interval.flush(4*60*1000+10);
      $httpBackend.flush();
    }
    expect($window.sessionStorage.token).toBe(token.toString());
    tokenRefresher.stop();
  });
  
  it('will not start if another instance is already running', function() {
    tokenRefresher.start();
    expect(tokenRefresher.isRunning()).toBe(true);
    tokenRefresher.stop();
    expect(tokenRefresher.isRunning()).toBe(false);
  });
  
  it('can manually update the token', function() {
    $httpBackend.expectGET(refreshTokenURL);
    tokenRefresher.refresh();
    $httpBackend.flush();
    expect($window.sessionStorage.token).toBe(token.toString());
  });
  
  it('can manually update the token multiple (12) times', function() {
    for (i = 0; i < 12; i++) {
      $httpBackend.expectGET(refreshTokenURL);
      tokenRefresher.refresh();
      $httpBackend.flush();
    }
    expect($window.sessionStorage.token).toBe(token.toString());
  });
  
  it('can stop updating token', function() {
    $httpBackend.expectGET(refreshTokenURL);
    tokenRefresher.start();
    $interval.flush(4*60*1000+10);
    $httpBackend.flush();
    expect($window.sessionStorage.token).toBe(token.toString());
    
    tokenRefresher.stop();
    $interval.flush(4*60*1000+10);
    expect($httpBackend.flush).toThrow(); // flush() should throw exception as there should be no pending request
  });
  
  it('should not start if no token is set', function() {
    // clear sessionstorage
    $window.sessionStorage.removeItem('token');
    
    tokenRefresher.start();
    $interval.flush(4*60*1000+10);
    
    expect($httpBackend.flush).toThrow(); // flush() should throw exception as there should be no pending request
  });
  
  /*
  it('dosn\'t freak out when something goes wrong', function() {
    expect(true).toBe(true);
  });
  */
});