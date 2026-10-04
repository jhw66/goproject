function addRow(){
    console.log("添加行");
    var table=document.getElementById("table");
    var length=table.rows.length;
    var newRow=table.insertRow(length);
    newRow.insertCell(0).innerHTML=length;
    newRow.insertCell(1).innerHTML="新添加的行";
    newRow.insertCell(2).innerHTML="新添加的年龄";
    newRow.insertCell(3).innerHTML='<button onclick="editRow(this)">编辑</button> <button onclick="deleteRow(this)">删除</button>';
}
function deleteRow(btn){
    console.log("删除行");
    var row=btn.parentNode.parentNode;//获取该节点的父节点的父节点，即tr节点
    row.parentNode.removeChild(row);
}
function editRow(btn){
    console.log("编辑行");
    var row=btn.parentNode.parentNode;
    var name=row.cells[1];
    var age=row.cells[2];
    var inputName=prompt("请输入新的名字:");
    var inputAge=prompt("请输入新的年龄:");
    if(inputName!=null && inputAge!=null){
        name.innerHTML=inputName;
        age.innerHTML=inputAge;
    }
}