toastr.options = {
    "closeButton": false,
    "debug": false,
    "newestOnTop": false,
    "progressBar": true,
    "positionClass": "toast-bottom-right",
    "preventDuplicates": true,
    "onclick": null,
    "showDuration": "300",
    "hideDuration": "1000",
    "timeOut": "5000",
    "extendedTimeOut": "1000",
    "showEasing": "swing",
    "hideEasing": "linear",
    "showMethod": "fadeIn",
    "hideMethod": "fadeOut"
};

let menusToApi = {
    'Roles':''
}

$(document).ready(function () {

    createSunKeyGraph("kube-system");
    $.getJSON("/namespaces", function (data) {
        let namespacesForSelect = [];
        data.forEach((elem, idx) => {
            namespacesForSelect.push({id: elem, text: elem})
        });
        namespacesForSelect.push({id: "all", text: "all"});
        let namespaceDropdown = $("#namespaces-dropdown");
        namespaceDropdown.select2({
            data: namespacesForSelect
        });
        namespaceDropdown.val("kube-system");
        namespaceDropdown.trigger('change');
    });
    $("#namespaces-dropdown").on("select2:select", function (e) {
        createSunKeyGraph(e.params.data.text);
    });

    $(".nav-item").on("click", function (e) {
        $(".nav-item").each((idx, el) => {
            $(el).removeClass("active")
        });
        $(e.currentTarget).addClass("active");
        console.log($(e.currentTarget).text().trim());

    });
});

function createSunKeyGraph(namespace) {

    $("#chart").empty();
    if (namespace === "all") namespace = "";


    d3.json("/roles?ns=" + namespace, function (error, json) {
        console.log(json);
        if (json.links === null || json.nodes.length === 0) {
            toastr.warning("Selected namespace: " + namespace + " doesn't have any roles");
        } else {
            var chart = d3.select("#chart").append("svg").chart("Sankey.Path");
            chart
                .colorNodes(function (name, node) {
                    return color(node, 1) || colors.fallback;
                })
                .colorLinks(function (link) {
                    return color(link.source, 4) || color(link.target, 1) || colors.fallback;
                })
                .nodeWidth(15)
                .nodePadding(10)
                .spread(true)
                .iterations(0)
                .draw(json);

            function color(node, depth) {
                return '#367d85';
            }
        }
    });

}