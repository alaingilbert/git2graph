var app = angular.module('app');

app.directive('project', function() {
  return {
    restrict: 'E',
    //scope: true,
    scope: {
      createDependency: '=',
      tree: '=',
      selectedNode: '=',
      colors: '=',
      dependency: '=',
      show: '='
    },
    template: '<div class="mysvg">' +
              '  <button class="btn btn-default" ng-click="setMode(\'commits\')"><i class="fa fa-circle"></i></button>' +
              '  <button class="btn btn-default" ng-click="setMode(\'links\')"><i class="fa fa-share-alt"></i></button>' +
              '  <button class="btn btn-default" ng-click="btnAddRowClicked()"><i class="fa fa-plus"></i> Add row</button>' +
              '  <button class="btn btn-default" ng-click="btnRedrawClicked()">Redraw</button>' +
              '  <svg></svg>' +
              '</div>',
    link: function(scope, element, attrs) {
      var $scope = scope;
      $scope.mode = 'commits';

      var firstCommit = null;
      var selectedPath = null;
      $scope.selectedNode = null;

      var lineFunction = d3.svg.line()
        .x(function(d) { return d.x; })
        .y(function(d) { return d.y; })
        .interpolate("linear");

      $scope.setMode = function(mode) {
        $scope.mode = mode;
        $scope.drawTree();
      };

      $scope.btnRedrawClicked = function() {
        $scope.drawTree();
      };

      $scope.btnAddRowClicked = function() {
        $scope.tree.push({id: $scope.tree.length.toString(), parents: [], column: 0, parentsPaths: {}, color: '#5aa1be'});
        $scope.drawTree();
      };

      $scope.drawTree = function() {
        if ($scope.$root.$$phase != '$apply' && $scope.$root.$$phase != '$digest') {
          $scope.$apply();
        }
        var xGap = 30;
        var yGap = 30;
        var cols = 10;
        var radius = 8;

        var svg = d3.select(element.find('svg')[0])
          .style('border', '1px solid gray')
          .attr('width', '100%')
          .attr('height', '300px');
        svg.selectAll('*').remove();

        var linesGroup = svg.append('g');
        var lineGroup = linesGroup.selectAll('lines')
          .data($scope.tree)
          .enter()
          .append('g');

        lineGroup
          .append('path')
          .attr('d', function(path, idx) {
            var d = [];
            d.push({x: 0, y: idx * yGap + yGap});
            d.push({x: $(svg[0][0]).width(), y: idx * yGap + yGap});
            return lineFunction(d)
          })
          .attr('stroke-width', 1)
          .attr('fill', 'none')
          .attr('stroke', '#aaa')
          .on('click', function(a, row, c) {
            $scope.tree.splice(row, 1);
            $scope.drawTree();
          });

          if ($scope.mode == 'commits') {
            lineGroup.selectAll('dummyCircle')
              .data(function(d, i) { return _.range(cols); })
              .enter()
              .append('circle')
                .attr('r', radius)
                .attr('stroke', '#ddd')
                .attr('fill', '#fff')
                .attr('cx', function(c) { return c * xGap + xGap; })
                .attr('cy', function(c, i, a) { return a * yGap + yGap })
                .on('mouseenter', function() { d3.select(this).attr('fill', '#f00'); })
                .on('mouseleave', function() { d3.select(this).attr('fill', '#fff'); })
                .on('click', function(a,b,c) {
                  $scope.tree[c].column = a;
                  $scope.drawTree();
                });
          } else {
            lineGroup.selectAll('dummyCircle')
              .data(function(d, i) { return _.range(cols); })
              .enter()
              .append('circle')
                .attr('r', 5)
                .attr('fill', 'none')
                .attr('stroke', '#ddd')
                .attr('fill', '#ddd')
                .attr('cx', function(c) { return c * xGap + xGap; })
                .attr('cy', function(c, i, a) { return a * yGap + yGap })
                .on('mousedown', function() {})
                .on('mouseup', function(item, columnIndex, rowIndex) {
                  if (selectedPath) {
                    console.log(selectedPath, columnIndex, rowIndex);
                    addPathNode(selectedPath[1], columnIndex, rowIndex);
                    console.log(selectedPath[1]);
                    selectedPath = null;
                    $scope.drawTree();
                  }
                });
          }

          var addPathNode = function(path, x, y) {
            if (y == path.path[0][1] && x > path.path[0][0]) {
              path.path.splice(1, 0, [x, y, 2]);
              return;
            }
            for (var i = 0; i < path.path.length; i++) {
              var pathItem = path.path[i];
              if (y > pathItem[1]) {
                continue;
              }
              path.path.splice(i, 0, [x, y, 1]);
              break;
            }
          };

          var nodesGroup = svg.append('g').attr('class', 'nodes');
          var nodeGroup = nodesGroup.selectAll('nodes')
            .data($scope.tree)
            .enter()
            .append('g')
            .attr('class', 'node');

          // Each paths
          var pathsGroup = nodeGroup
            .append('g')
            .attr('class', 'paths');
          pathsGroup
            .selectAll('g')
            .data(function(node) {
              return _.map(node.parentsPaths, function(pathItem, parentId) { return [node, pathItem, parentId]; });
            })
            .enter()
            .append('path')
              .attr('d', function(item) {
                var node = item[0];
                var path = item[1];
                var parentId = item[2];
                var p = $scope.tree[parentId];
                var result = [];
                _.each(path.path, function(item) {
                  var point = {
                    x: item[0] * xGap + xGap,
                    y: item[1] * yGap + yGap
                  };
                  if (item[2] == 2) {
                    point.y += 6;
                  }
                  if (item[2] == 1) {
                    point.y -= 6;
                  }
                  result.push(point);
                });
                return lineFunction(result)
              })
              .attr('stroke-width', 3)
              .attr('fill', 'none')
              .attr('stroke', function(item) {
                return item[1].color;
              })
              .on('mousedown', function(item) {
                selectedPath = item;
              })
              .on('click', function(item) {
                var node = item[0];
                var idx = _.indexOf(node.parents, item[1]);
                $scope.tree[node.id].parents.splice(idx, 1);
                $scope.drawTree();
              });

          var addPath = function(rowIndex) {
            var firstNode = $scope.tree[firstCommit];
            var lastNode = $scope.tree[rowIndex];
            $scope.tree[firstCommit].parents.push(rowIndex.toString());
            $scope.tree[firstCommit].parents = _.sortBy($scope.tree[firstCommit].parents);
            $scope.tree[firstCommit].parentsPaths[rowIndex.toString()] = {
              path: [
                [firstNode.column, parseInt(firstNode.id), 0],
                [lastNode.column, parseInt(lastNode.id), 0]],
              color: '#5aa1be'};
            firstCommit = null;
            $scope.drawTree();
          };

          // Each nodes
          nodeGroup.append('circle')
            .attr('r', radius)
            .attr('fill', function(node) { return node.color; })
            .attr('stroke', '#000')
            .attr('cx', function(node) {
              return node.column * xGap + xGap;
            })
            .attr('cy', function(node, idx) { return idx * yGap + yGap })
            .on('mouseenter', function() { d3.select(this).attr('fill', '#f00'); })
            .on('mouseleave', function(node) { d3.select(this).attr('fill', node.color); })
            .on('mousedown', function(node, rowIndex) {
              firstCommit = rowIndex;
            })
            .on('mouseup', function(node, rowIndex) {
              if (firstCommit == rowIndex) {
                return;
              }
              if (firstCommit > rowIndex) {
                firstCommit = [rowIndex, rowIndex = firstCommit][0]; // swap vars
              }
              if (!_.includes($scope.tree[firstCommit].parents, rowIndex)) {
                addPath(rowIndex);
              }
            })
            .on('click', function(node) {
              $scope.selectedNode = node;
              $scope.$apply();
            });

      };

      $scope.$watchCollection('tree', function(newValue, oldValue) {
        $scope.drawTree();
        //new ZeroClipboard($('.copy-btn'));
      });

    }
  }
});
