/**
 * backgroundImage makes sure a background image can be set
 * from html and also be two-way binded to a variable
 */
app.directive('backgroundImage', function() {
  return {
    scope: {
      path: '@'
    },
    link: function(scope, element, attrs) {
      scope.$parent.$watch(scope.path, function() {
        var url = attrs.backgroundImage;
        
        element.css({
            'background-image': 'url("' + url +'")',
            'background-size' : 'cover'
        });
      });
    }
  };
});