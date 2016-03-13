var app = angular.module('app');

app.directive('project', function() {
  return {
    restrict: 'E',
    //scope: true,
    scope: {
      createDependency: '=',
      inputFile: '=',
      testFile: '=',
      shellFile: '=',
      tree: '=',
      selectedNode: '=',
      colors: '=',
      dependency: '=',
      show: '='
    },
    template: '<div class="mysvg">' +
              '  <button class="btn btn-default" ng-click="btnAddRowClicked()"><i class="fa fa-plus"></i> Add row</button>' +
              '  <button class="btn btn-default" ng-click="btnRedrawClicked()">Redraw</button>' +
              '  <svg></svg>' +
              '</div>',
    link: function(scope, element, attrs) {
      var $scope = scope;

      var firstCommit = null;
      var selectedPath = null;
      $scope.selectedNode = null;

      var lineFunction = d3.svg.line()
        .x(function(d) { return d.x; })
        .y(function(d) { return d.y; })
        .interpolate("linear");

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
          .style('height', ($scope.tree.length + 1) * yGap + radius);
        svg.selectAll('*').remove();

        var linesGroup = svg.append('g');
        var lineGroup = linesGroup.selectAll('lines')
          .data($scope.tree)
          .enter()
          .append('g');

        lineGroup
          .append('text')
          .attr('x', 3)
          .attr('y', function(path, idx) {
            return yGap * idx + yGap;
          })
          .attr('alignment-baseline', 'middle')
          .attr('text-anchor', 'left')
          .text(function(path, idx) { return idx; });

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

          lineGroup.selectAll('dummyCircle')
            .data(function(d, i) { return _.range(cols); })
            .enter()
            .append('circle')
              .attr('r', 5)
              .attr('stroke', '#ddd')
              .attr('fill', '#ddd')
              .attr('cx', function(c) { return c * xGap + xGap; })
              .attr('cy', function(c, i, a) { return a * yGap + yGap })
              .on('mouseenter', function() { d3.select(this).attr('fill', '#aaa'); })
              .on('mouseleave', function() { d3.select(this).attr('fill', '#ddd'); })
              .on('mousedown', function() {})
              .on('mouseup', function(item, columnIndex, rowIndex) {
                if (selectedPath) {
                  addPathNode(selectedPath[1], columnIndex, rowIndex);
                  selectedPath = null;
                  $scope.drawTree();
                }
              })
              .on('click', function(a,b,c) {
                $scope.tree[c].column = a;
                $scope.drawTree();
              });

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
              if (y == pathItem[1] && x < pathItem[0]) {
                path.path.splice(i+1, 0, [x, y, 0]);
              } else {
                path.path.splice(i, 0, [x, y, 1]);
              }
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
              .on('mouseenter', function() {
                d3.select(this).attr('opacity', 0.4);
              })
              .on('mouseleave', function(node) {
                d3.select(this).attr('opacity', 1);
              })
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
            .attr('opacity', function(node) {
              if ($scope.selectedNode && $scope.selectedNode.id == node.id) {
                return 0.4;
              } else {
                return 1;
              }
            })
            .attr('stroke', '#000')
            .attr('cx', function(node) {
              return node.column * xGap + xGap;
            })
            .attr('cy', function(node, idx) { return idx * yGap + yGap })
            .on('mouseenter', function() {
              d3.select(this).attr('opacity', 0.4);
            })
            .on('mouseleave', function(node) {
              if ($scope.selectedNode && $scope.selectedNode.id == node.id) {
                d3.select(this).attr('opacity', 0.4);
              } else {
                d3.select(this).attr('opacity', 1);
              }
            })
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
              if ($scope.selectedNode && $scope.selectedNode.id == node.id) {
                $scope.selectedNode = null;
              } else {
                $scope.selectedNode = node;
              }
              $scope.drawTree();
            });

        generateInputFile();
        generateTestsFile();
        generateShellScript();
      };


      var generateInputFile = function() {
        var output = _.map($scope.tree, function(el) { return _.pick(el, ['id', 'parents']); });
        $scope.inputFile = JSON.stringify(output, null, 2);
      };


      var generateTestsFile = function() {
        var testName = 'Test1';
        var out = '';
        out += 'func ' + testName + '(t *testing.T) {\n';

        out += '	// Initial input\n';
        out += '	inputNodes := make([]map[string]interface{}, 0)\n';
        _.each($scope.tree, function(node) {
          var p = _.map(node.parents, function(el) { return '"' + el + '"'; }).join(',');
          out += '	inputNodes = append(inputNodes, map[string]interface{}{"id": "' + node.id + '", "parents": []string{' + p + '}})\n';
        });

        out += '\n	out, _ := buildTree(inputNodes, customColors)\n\n';

        out += '	// Expected output\n';
        var expectedColumns = _.map($scope.tree, 'column').join(', ');
        out += '	expectedColumns := []int{' + expectedColumns + '}\n\n';

        out += '	expectedPaths := []map[string]Path{\n';
        _.each($scope.tree, function(node) {
          out += '		map[string]Path{\n';
          _.each(node.parents, function(parentId) {
            var parentNode = $scope.tree[parentId];
            out += '			"' + parentId + '": Path{"' + parentId + '", []Point{Point{' + node.column + ', ' + node.id + ', 0}, Point{' + parentNode.column + ', ' + parentNode.id + ', 0}}, "nocolor"},\n';
          });
          out += '		},\n';
        });
        out += '	}\n\n';

        out += '	// Validation\n';
        out += '	validateColumns(t, expectedColumns, out)\n';
        out += '	validatePaths(t, expectedPaths, out)\n';
        out += '	validateColors(t, expectedPaths, out)\n';
        out += '}';

        $scope.testFile = out;
      };


      var generateShellScript = function() {

        var createCommit = function(id) {
          var out = '';
          out += 'touch ' + id + '\n';
          out += 'git add ' + id + '\n';
          out += 'git commit -m ' + id + '\n';
          return out;
        };

        var reversedNodes = _.cloneDeep($scope.tree);
        _.reverse(reversedNodes);
        var out = "";
        _.each(reversedNodes, function(item, idx) {
          item.parents = _.keys(item.parentsPaths);
          item.parents = _.sortBy(item.parents, function(item) { return $scope.tree[item].column; });
          if (item.parents.length == 0) {
            out += 'git checkout -b ' + item.id + '\n';
            out += createCommit(item.id);
          } else if (item.parents.length == 1) {
            if ($scope.tree[item.parents[0]].column < item.column) {
              out += 'git checkout ' + item.parents[0] + '\n';
              out += 'git checkout -b ' + item.id + '\n';
              out += createCommit(item.id);
            } else if ($scope.tree[item.parents[0]].column == item.column) {
              out += 'git checkout ' + item.parents[0] + '\n';
              out += 'git checkout -b ' + item.id + '\n';
              out += createCommit(item.id);
            }
          } else if (item.parents.length == 2) {
            if ($scope.tree[item.parents[0]].column < item.column) {
              out += 'git checkout ' + item.parents[1] + '\n';
              out += 'git checkout -b ' + item.id + '\n';
              out += 'git merge -m ' + item.parents[0] + ' --no-ff ' + item.parents[0] + '\n';
            } else if ($scope.tree[item.parents[1]].column > item.column) {
              out += 'git checkout ' + item.parents[0] + '\n';
              out += 'git checkout -b ' + item.id + '\n';
              out += 'git merge -m ' + item.parents[1] + ' --no-ff ' + item.parents[1] + '\n';
            }
          }
        });

        $scope.shellFile = out;
      };


      $scope.$watch('tree', function(newValue, oldValue) {
        $scope.drawTree();
      }, true);

    }
  }
});


app.directive('focusMe', function($timeout) {
  return {
    scope: { trigger: '@focusMe' },
    link: function(scope, element, attr) {
      var predicate = attr.focusMe || true;
      switch (predicate) {
        case "true": case "yes": case "1":
          predicate = true;
          break;
        case "false": case "no": case "0": case null:
          predicate = false;
          break;
        default:
          predicate = Boolean(predicate);
          break;
      }
      if (!predicate) {
        return;
      }
      scope.$watch('trigger', function(value) {
        $timeout(function() {
          element[0].focus();
        });
      });
    }
  };
});


app.directive('ngEnter', function() {
  return function(scope, element, attrs) {
    element.bind("keydown keypress", function(event) {
      if(event.which === 13) {
        scope.$apply(function(){
          scope.$eval(attrs.ngEnter, {'event': event});
        });
        event.preventDefault();
      }
    });
  };
});
