$('#first-datatable-markup pre code').text($('#first-datatable-without').html());
$('#first-datatable-output').html($('#first-datatable-without').html());
var firstExampleTable = document.querySelector('#first-datatable-output table');
var datatable = new DataTable(firstExampleTable, {
    pageSize: 10,
    pagingDivSelector: '#paging-first-datatable',
    sort: [true, true, true, true],
    filters: [false, true, 'select', true],
    filterText: 'Filtrar... '
});
function reload() {
    window.location.reload();
}