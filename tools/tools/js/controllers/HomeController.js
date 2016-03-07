var app = angular.module('app', ['ngSanitize', 'ui.bootstrap', 'LocalStorageModule',]);

app.controller('HomeController',
  function($scope, $sce)
  {

    $scope.btnSaveClicked = function() {
      localStorage.setItem('tree', JSON.stringify($scope.tree));
    };

    $scope.btnLoadClicked = function() {
      var savedTree = localStorage.getItem('tree');
      $scope.tree = JSON.parse(savedTree);
    };


    var highlightJson = function(json) {
      json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
      return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function(match) {
        var cls = 'number';
        if (/^"/.test(match)) {
          if (/:$/.test(match)) {
            cls = 'key';
          } else {
            cls = 'string';
          }
        } else if (/true|false/.test(match)) {
          cls = 'boolean';
        } else if (/null/.test(match)) {
          cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
      });
    };

    $scope.highlight = function(obj) {
      return highlightJson(angular.toJson(obj, true))
    };


    (function constructor() {
      $scope.selectedNode = null;
      $scope.tree = [];

      $scope.colors = [
        '#5aa1be',
        '#c065b8',
        '#c0ab5f',
        '#59bc95',
        '#7a63be',
        '#c0615b',
        '#73bb5e',
        '#6ee585',
        '#7088e8',
        '#eb77a3',
        '#c2e675',
      ];

    })();

  }
);
