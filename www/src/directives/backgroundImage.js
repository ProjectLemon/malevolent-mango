app.directive('backgroundImage', function() {
  return {
    scope: {
      path: '@'
    },
    link: function(scope, element, attrs) {
      scope.$parent.$watch(scope.path, function() {
        var url = attrs.backgroundImage;
        
        element.css({
            'background-image': 'url(' + url +')',
            'background-size' : 'cover'
        });
      });
    }
  };
});