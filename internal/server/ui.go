package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Bounty</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#c45d2c;--rl:#e8753a;--leather:#a0845c;--ll:#c4a87a;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c44040;--blue:#4a7ec4;--mono:'JetBrains Mono',Consolas,monospace;--serif:'Libre Baskerville',Georgia,serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);font-size:13px;line-height:1.6}
a{color:var(--rl);text-decoration:none}a:hover{color:var(--gold)}
.hdr{padding:.6rem 1.2rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem}
.hdr h1{font-family:var(--serif);font-size:1rem;white-space:nowrap}.hdr h1 span{color:var(--rl)}
.hdr-right{display:flex;align-items:center;gap:1rem;font-size:.7rem;color:var(--leather)}
.hdr-right b{color:var(--cream)}
.main{max-width:900px;margin:0 auto;padding:1rem 1.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.toolbar select,.toolbar input[type=text]{background:var(--bg);border:1px solid var(--bg3);color:var(--cream);padding:.35rem .5rem;font-family:var(--mono);font-size:.75rem;outline:none}
.toolbar select:focus,.toolbar input:focus{border-color:var(--rust)}
.toolbar input[type=text]{flex:1;min-width:120px}
.btn{font-family:var(--mono);font-size:.72rem;padding:.35rem .7rem;border:1px solid;cursor:pointer;background:transparent;transition:all .15s;white-space:nowrap}
.btn-p{border-color:var(--rust);color:var(--rl)}.btn-p:hover{background:var(--rust);color:var(--cream)}
.btn-g{border-color:var(--gold);color:var(--gold)}.btn-g:hover{background:var(--gold);color:var(--bg)}
.btn-d{border-color:var(--bg3);color:var(--cm)}.btn-d:hover{border-color:var(--red);color:var(--red)}
.btn-s{border-color:var(--green);color:var(--green)}.btn-s:hover{background:var(--green);color:var(--bg)}
.btn-b{border-color:var(--blue);color:var(--blue)}.btn-b:hover{background:var(--blue);color:var(--cream)}
.tabs{display:flex;gap:0;margin-bottom:1rem;border-bottom:1px solid var(--bg3)}
.tab{padding:.4rem 1rem;cursor:pointer;font-size:.75rem;color:var(--cm);border-bottom:2px solid transparent;transition:.15s}
.tab:hover{color:var(--cream)}.tab.active{color:var(--rl);border-bottom-color:var(--rl)}

.issue-row{display:flex;align-items:flex-start;gap:.6rem;padding:.55rem .7rem;border:1px solid var(--bg3);background:var(--bg2);margin-bottom:1px;cursor:pointer;transition:background .1s}
.issue-row:hover{background:var(--bg3)}
.issue-status{width:10px;height:10px;border-radius:50%;margin-top:5px;flex-shrink:0}
.issue-status.open{background:var(--green)}.issue-status.in_progress{background:var(--gold)}.issue-status.closed{background:var(--cm)}
.issue-body{flex:1;min-width:0}
.issue-title{font-size:.8rem;font-weight:600;color:var(--cream)}
.issue-meta{font-size:.65rem;color:var(--cm);margin-top:2px;display:flex;gap:.7rem;flex-wrap:wrap}
.issue-num{color:var(--leather)}
.priority-badge{font-size:.6rem;padding:0 .35rem;border:1px solid;border-radius:2px;text-transform:uppercase;letter-spacing:.5px}
.pri-critical{border-color:var(--red);color:var(--red)}.pri-high{border-color:var(--rl);color:var(--rl)}.pri-medium{border-color:var(--gold);color:var(--gold)}.pri-low{border-color:var(--cm);color:var(--cm)}
.label-chip{font-size:.6rem;padding:0 .3rem;background:var(--bg3);color:var(--ll);border-radius:2px}
.comment-count{color:var(--leather)}

.modal-bg{position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.65);display:flex;align-items:center;justify-content:center;z-index:100}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:95%;max-width:700px;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--serif);font-size:.95rem;margin-bottom:1rem}
label.fl{display:block;font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px;margin-bottom:.25rem;margin-top:.7rem}
input[type=text],input[type=date],textarea,select{background:var(--bg);border:1px solid var(--bg3);color:var(--cream);padding:.4rem .6rem;font-family:var(--mono);font-size:.8rem;width:100%;outline:none}
input:focus,textarea:focus,select:focus{border-color:var(--rust)}
textarea{resize:vertical;min-height:80px}
.form-row{display:flex;gap:.5rem}.form-row>*{flex:1}
.comment-thread{margin-top:1rem;border-top:1px solid var(--bg3);padding-top:.8rem}
.comment-item{padding:.5rem 0;border-bottom:1px solid var(--bg3)}
.comment-item:last-child{border-bottom:none}
.comment-author{font-size:.7rem;color:var(--rl);font-weight:600}
.comment-time{font-size:.6rem;color:var(--cm);margin-left:.5rem}
.comment-body{font-size:.78rem;color:var(--cd);margin-top:.2rem;white-space:pre-wrap}
.empty{text-align:center;padding:2rem;color:var(--cm);font-style:italic;font-family:var(--serif)}

.ms-card{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem;margin-bottom:.5rem}
.ms-title{font-size:.8rem;font-weight:600}.ms-meta{font-size:.65rem;color:var(--cm);margin-top:.2rem}
.ms-bar{height:4px;background:var(--bg3);margin-top:.4rem;border-radius:2px;overflow:hidden}
.ms-fill{height:100%;background:var(--green);transition:width .3s}

.act-item{font-size:.72rem;padding:.3rem 0;border-bottom:1px solid var(--bg3);color:var(--cd)}
.act-action{font-weight:600;color:var(--rl)}.act-time{color:var(--cm);font-size:.6rem;float:right}

.board{display:grid;grid-template-columns:repeat(3,1fr);gap:.6rem;min-height:400px}
@media(max-width:700px){.board{grid-template-columns:1fr}}
.board-col{background:var(--bg2);border:1px solid var(--bg3);padding:.5rem;min-height:300px}
.board-col.drag-over{border-color:var(--rust);background:rgba(196,93,44,.08)}
.col-hdr{font-size:.7rem;text-transform:uppercase;letter-spacing:1px;color:var(--leather);padding:.3rem .2rem .5rem;display:flex;justify-content:space-between;align-items:center}
.col-hdr .col-count{background:var(--bg3);color:var(--cream);padding:.1rem .4rem;border-radius:2px;font-size:.6rem}
.col-hdr .col-open{color:var(--green)}.col-hdr .col-wip{color:var(--gold)}.col-hdr .col-done{color:var(--cm)}
.board-card{background:var(--bg);border:1px solid var(--bg3);padding:.5rem .6rem;margin-bottom:.35rem;cursor:grab;transition:all .15s;font-size:.75rem}
.board-card:hover{border-color:var(--leather)}
.board-card:active{cursor:grabbing}
.board-card.dragging{opacity:.4;border-color:var(--rust)}
.board-card .bc-title{font-weight:600;color:var(--cream);margin-bottom:.2rem}
.board-card .bc-meta{font-size:.6rem;color:var(--cm);display:flex;gap:.5rem;flex-wrap:wrap}
.board-card .bc-num{color:var(--leather)}
.view-toggle{display:flex;gap:0;margin-left:auto}
.view-toggle .vt{padding:.25rem .6rem;font-size:.65rem;cursor:pointer;border:1px solid var(--bg3);color:var(--cm);background:var(--bg)}
.view-toggle .vt:first-child{border-right:none}
.view-toggle .vt.active{background:var(--bg2);color:var(--rl);border-color:var(--rust)}

.proj-card{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem;margin-bottom:.5rem;cursor:pointer;transition:background .1s}
.proj-card:hover{background:var(--bg3)}
.proj-card h3{font-size:.85rem;margin-bottom:.2rem}
.proj-card .proj-stats{font-size:.65rem;color:var(--cm)}
.proj-card .proj-stats b{color:var(--green)}
</style>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital@0;1&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
</head><body>
<div class="hdr">
<h1><span>Bounty</span></h1>
<div class="hdr-right">
<select id="projSelect" onchange="switchProject(this.value)" style="background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.72rem;padding:.25rem .4rem"></select>
<span>Open: <b id="sOpen">-</b></span>
<span>Closed: <b id="sClosed">-</b></span>
<button class="btn btn-p" onclick="showNewProject()">+ Project</button>
</div>
</div>
<div class="main">
<div id="upgrade-banner" class="upgrade" style="display:none">
  <strong style="color:var(--cream)">Free tier</strong> — 2 projects, 50 issues. <a href="https://stockyard.dev/bounty/" target="_blank">Upgrade to Pro for $0.99/mo →</a>
</div>

<div class="tabs">
<div class="tab active" data-tab="issues" onclick="switchTab('issues')">Issues</div>
<div class="tab" data-tab="board" onclick="switchTab('board')">Board</div>
<div class="tab" data-tab="milestones" onclick="switchTab('milestones')">Milestones</div>
<div class="tab" data-tab="activity" onclick="switchTab('activity')">Activity</div>
</div>

<div id="pane-issues">
<div class="toolbar">
<select id="fStatus" onchange="loadIssues()"><option value="open">Open</option><option value="in_progress">In Progress</option><option value="closed">Closed</option><option value="all">All</option></select>
<select id="fPriority" onchange="loadIssues()"><option value="">Priority</option><option value="critical">Critical</option><option value="high">High</option><option value="medium">Medium</option><option value="low">Low</option></select>
<select id="fLabel" onchange="loadIssues()"><option value="">Label</option></select>
<select id="fAssignee" onchange="loadIssues()"><option value="">Assignee</option></select>
<select id="fSort" onchange="loadIssues()"><option value="created">Newest</option><option value="updated">Updated</option><option value="priority">Priority</option><option value="comments">Comments</option></select>
<input type="text" id="fSearch" placeholder="Search issues..." onkeydown="if(event.key==='Enter')loadIssues()">
<button class="btn btn-p" onclick="showNewIssue()">+ Issue</button>
</div>
<div id="issueList"></div>
</div>

<div id="pane-board" style="display:none">
<div class="board">
<div class="board-col" data-status="open" ondragover="boardDragOver(event)" ondragleave="boardDragLeave(event)" ondrop="boardDrop(event)">
<div class="col-hdr"><span class="col-open">Open</span><span class="col-count" id="bc-open">0</span></div>
<div class="col-cards" id="col-open"></div>
</div>
<div class="board-col" data-status="in_progress" ondragover="boardDragOver(event)" ondragleave="boardDragLeave(event)" ondrop="boardDrop(event)">
<div class="col-hdr"><span class="col-wip">In Progress</span><span class="col-count" id="bc-wip">0</span></div>
<div class="col-cards" id="col-in_progress"></div>
</div>
<div class="board-col" data-status="closed" ondragover="boardDragOver(event)" ondragleave="boardDragLeave(event)" ondrop="boardDrop(event)">
<div class="col-hdr"><span class="col-done">Closed</span><span class="col-count" id="bc-closed">0</span></div>
<div class="col-cards" id="col-closed"></div>
</div>
</div>
</div>

<div id="pane-milestones" style="display:none">
<div style="display:flex;justify-content:space-between;margin-bottom:.8rem">
<span style="font-size:.7rem;color:var(--leather)">Project milestones</span>
<button class="btn btn-p" onclick="showNewMilestone()">+ Milestone</button>
</div>
<div id="msList"></div>
</div>

<div id="pane-activity" style="display:none">
<div id="actList"></div>
</div>
</div>
<div id="modal"></div>

<script>
let projects=[],curProject='',issues=[],milestones=[];

async function api(url,opts){const r=await fetch(url,opts);return r.json()}
function esc(s){return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;')}
function timeAgo(d){const s=Math.floor((Date.now()-new Date(d))/1e3);if(s<60)return s+'s ago';if(s<3600)return Math.floor(s/60)+'m ago';if(s<86400)return Math.floor(s/3600)+'h ago';return Math.floor(s/86400)+'d ago'}

async function init(){
  const d=await api('/api/projects');projects=d.projects||[];
  const sel=document.getElementById('projSelect');
  sel.innerHTML='<option value="">All projects</option>'+projects.map(p=>'<option value="'+p.id+'">'+esc(p.name)+'</option>').join('');
  if(curProject)sel.value=curProject;
  loadStats();loadIssues();loadFilters();
}

function switchProject(id){curProject=id;loadIssues();loadStats();loadMilestones();loadActivity();loadFilters()}

async function loadStats(){
  const d=await api('/api/stats');
  document.getElementById('sOpen').textContent=d.open_issues;
  document.getElementById('sClosed').textContent=d.closed_issues;
}

async function loadFilters(){
  const [lb,as]=await Promise.all([api('/api/labels'),api('/api/assignees')]);
  const ls=document.getElementById('fLabel');ls.innerHTML='<option value="">Label</option>'+(lb.labels||[]).map(l=>'<option>'+esc(l)+'</option>').join('');
  const aa=document.getElementById('fAssignee');aa.innerHTML='<option value="">Assignee</option>'+(as.assignees||[]).map(a=>'<option>'+esc(a)+'</option>').join('');
}

async function loadIssues(){
  const p=new URLSearchParams();
  if(curProject)p.set('project_id',curProject);
  p.set('status',document.getElementById('fStatus').value);
  const pri=document.getElementById('fPriority').value;if(pri)p.set('priority',pri);
  const lab=document.getElementById('fLabel').value;if(lab)p.set('label',lab);
  const asn=document.getElementById('fAssignee').value;if(asn)p.set('assignee',asn);
  const srt=document.getElementById('fSort').value;p.set('sort',srt);
  const srch=document.getElementById('fSearch').value;if(srch)p.set('search',srch);
  const d=await api('/api/issues?'+p);issues=d.issues||[];
  renderIssues();
}

function renderIssues(){
  const el=document.getElementById('issueList');
  if(!issues.length){el.innerHTML='<div class="empty">No issues found.</div>';return}
  el.innerHTML=issues.map(i=>{
    const priCls='pri-'+i.priority;
    const labels=(i.labels||[]).map(l=>'<span class="label-chip">'+esc(l)+'</span>').join(' ');
    const comments=i.comment_count?'<span class="comment-count">&#x1f4ac; '+i.comment_count+'</span>':'';
    const assignee=i.assignee?'<span>&#x1f464; '+esc(i.assignee)+'</span>':'';
    return '<div class="issue-row" onclick="showIssue(\''+i.id+'\')">'+
      '<div class="issue-status '+i.status+'"></div>'+
      '<div class="issue-body">'+
        '<div class="issue-title">'+esc(i.title)+'</div>'+
        '<div class="issue-meta">'+
          '<span class="issue-num">#'+i.number+'</span>'+
          '<span class="priority-badge '+priCls+'">'+i.priority+'</span>'+
          labels+comments+assignee+
          '<span>'+timeAgo(i.created_at)+'</span>'+
        '</div>'+
      '</div></div>'
  }).join('')
}

async function showIssue(id){
  const [i,cd]=await Promise.all([api('/api/issues/'+id),api('/api/issues/'+id+'/comments')]);
  const comments=(cd.comments||[]).map(c=>'<div class="comment-item"><span class="comment-author">'+esc(c.author||'anonymous')+'</span><span class="comment-time">'+timeAgo(c.created_at)+'</span><div class="comment-body">'+esc(c.body)+'</div></div>').join('');
  const priCls='pri-'+i.priority;
  const labels=(i.labels||[]).map(l=>'<span class="label-chip">'+esc(l)+'</span>').join(' ');
  const statusBtn=i.status==='open'?
    '<button class="btn btn-b" onclick="setStatus(\''+i.id+'\',\'in_progress\')">Start</button><button class="btn btn-d" onclick="closeIss(\''+i.id+'\')">Close</button>':
    i.status==='in_progress'?
    '<button class="btn btn-s" onclick="reopenIss(\''+i.id+'\')">Back to Open</button><button class="btn btn-d" onclick="closeIss(\''+i.id+'\')">Close</button>':
    '<button class="btn btn-s" onclick="reopenIss(\''+i.id+'\')">Reopen</button>';
  const proj=projects.find(p=>p.id===i.project_id);
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<div style="display:flex;justify-content:space-between;align-items:flex-start">'+
      '<h2><span class="issue-num">#'+i.number+'</span> '+esc(i.title)+'</h2>'+
      '<div style="display:flex;gap:.3rem">'+
        '<button class="btn btn-g" onclick="showEditIssue(\''+i.id+'\')">Edit</button>'+
        statusBtn+
        '<button class="btn btn-d" onclick="if(confirm(\'Delete this issue?\'))delIss(\''+i.id+'\')">Del</button>'+
      '</div>'+
    '</div>'+
    '<div style="display:flex;gap:1rem;flex-wrap:wrap;margin:.6rem 0;font-size:.7rem;color:var(--leather)">'+
      '<span class="priority-badge '+priCls+'">'+i.priority+'</span>'+
      '<span>Status: <b style="color:'+(i.status==='open'?'var(--green)':i.status==='in_progress'?'var(--gold)':'var(--cm)')+'">'+i.status.replace('_',' ')+'</b></span>'+
      (proj?'<span>Project: '+esc(proj.name)+'</span>':'')+
      (i.assignee?'<span>Assignee: '+esc(i.assignee)+'</span>':'')+
      (i.time_estimate?'<span>Est: '+i.time_estimate+'m</span>':'')+
      (i.time_spent?'<span>Spent: '+i.time_spent+'m</span>':'')+
      ' '+labels+
    '</div>'+
    (i.body?'<div style="padding:.7rem;background:var(--bg);border:1px solid var(--bg3);font-size:.78rem;color:var(--cd);white-space:pre-wrap;margin:.5rem 0">'+esc(i.body)+'</div>':'')+
    '<div class="comment-thread">'+
      '<div style="font-size:.7rem;color:var(--leather);margin-bottom:.5rem">Comments ('+(cd.comments||[]).length+')</div>'+
      (comments||'<div style="font-size:.75rem;color:var(--cm);font-style:italic">No comments yet.</div>')+
      '<div style="margin-top:.8rem">'+
        '<input type="text" id="cmtAuthor" placeholder="Your name" style="margin-bottom:.3rem">'+
        '<textarea id="cmtBody" placeholder="Add a comment..." rows="3"></textarea>'+
        '<button class="btn btn-p" style="margin-top:.3rem" onclick="addComment(\''+i.id+'\')">Comment</button>'+
      '</div>'+
    '</div>'+
  '</div></div>'
}

async function addComment(issueID){
  const author=document.getElementById('cmtAuthor').value.trim()||'anonymous';
  const body=document.getElementById('cmtBody').value.trim();
  if(!body)return;
  await api('/api/issues/'+issueID+'/comments',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({author,body})});
  showIssue(issueID);loadStats();
}

async function closeIss(id){await api('/api/issues/'+id+'/close',{method:'POST'});closeModal();loadIssues();loadStats()}
async function reopenIss(id){await api('/api/issues/'+id+'/reopen',{method:'POST'});closeModal();loadIssues();loadStats()}
async function setStatus(id,status){await api('/api/issues/'+id+'/status',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({status})});closeModal();loadIssues();loadStats()}
async function delIss(id){await api('/api/issues/'+id,{method:'DELETE'});closeModal();loadIssues();loadStats()}

function showNewIssue(){
  if(!projects.length){alert('Create a project first');return}
  const projOpts=projects.map(p=>'<option value="'+p.id+'"'+(p.id===curProject?' selected':'')+'>'+esc(p.name)+'</option>').join('');
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>New Issue</h2>'+
    '<label class="fl">Project</label><select id="ni-proj">'+projOpts+'</select>'+
    '<label class="fl">Title</label><input type="text" id="ni-title">'+
    '<label class="fl">Description</label><textarea id="ni-body" rows="4"></textarea>'+
    '<div class="form-row">'+
      '<div><label class="fl">Priority</label><select id="ni-pri"><option value="low">Low</option><option value="medium" selected>Medium</option><option value="high">High</option><option value="critical">Critical</option></select></div>'+
      '<div><label class="fl">Assignee</label><input type="text" id="ni-assign"></div>'+
    '</div>'+
    '<label class="fl">Labels (comma-separated)</label><input type="text" id="ni-labels" placeholder="bug, frontend">'+
    '<div class="form-row">'+
      '<div><label class="fl">Time estimate (min)</label><input type="text" id="ni-est" placeholder="0"></div>'+
      '<div><label class="fl">Milestone</label><select id="ni-ms"><option value="">None</option></select></div>'+
    '</div>'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveNewIssue()">Create</button><button class="btn btn-d" onclick="closeModal()">Cancel</button></div>'+
  '</div></div>';
  loadMsSelect()
}

async function loadMsSelect(){
  const pid=document.getElementById('ni-proj')?document.getElementById('ni-proj').value:'';
  if(!pid)return;
  const d=await api('/api/projects/'+pid+'/milestones');
  const sel=document.getElementById('ni-ms');
  if(sel)sel.innerHTML='<option value="">None</option>'+((d.milestones||[]).map(m=>'<option value="'+m.id+'">'+esc(m.title)+'</option>').join(''));
}

async function saveNewIssue(){
  const labels=(document.getElementById('ni-labels').value||'').split(',').map(s=>s.trim()).filter(Boolean);
  const body={
    project_id:document.getElementById('ni-proj').value,
    title:document.getElementById('ni-title').value,
    body:document.getElementById('ni-body').value,
    priority:document.getElementById('ni-pri').value,
    assignee:document.getElementById('ni-assign').value,
    labels:labels,
    time_estimate:parseInt(document.getElementById('ni-est').value)||0,
    milestone_id:document.getElementById('ni-ms').value||''
  };
  if(!body.title){alert('Title required');return}
  const r=await api('/api/issues',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  if(r.error){alert(r.error);return}
  closeModal();loadIssues();loadStats();loadFilters();
}

function showEditIssue(id){
  const i=issues.find(x=>x.id===id);
  if(!i)return;
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>Edit #'+i.number+'</h2>'+
    '<label class="fl">Title</label><input type="text" id="ei-title" value="'+esc(i.title)+'">'+
    '<label class="fl">Description</label><textarea id="ei-body" rows="4">'+esc(i.body)+'</textarea>'+
    '<div class="form-row">'+
      '<div><label class="fl">Priority</label><select id="ei-pri"><option value="low"'+(i.priority==='low'?' selected':'')+'>Low</option><option value="medium"'+(i.priority==='medium'?' selected':'')+'>Medium</option><option value="high"'+(i.priority==='high'?' selected':'')+'>High</option><option value="critical"'+(i.priority==='critical'?' selected':'')+'>Critical</option></select></div>'+
      '<div><label class="fl">Assignee</label><input type="text" id="ei-assign" value="'+esc(i.assignee)+'"></div>'+
    '</div>'+
    '<label class="fl">Labels (comma-separated)</label><input type="text" id="ei-labels" value="'+esc((i.labels||[]).join(', '))+'">'+
    '<div class="form-row">'+
      '<div><label class="fl">Time estimate (min)</label><input type="text" id="ei-est" value="'+(i.time_estimate||'')+'"></div>'+
      '<div><label class="fl">Time spent (min)</label><input type="text" id="ei-spent" value="'+(i.time_spent||'')+'"></div>'+
    '</div>'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveEditIssue(\''+i.id+'\')">Save</button><button class="btn btn-d" onclick="showIssue(\''+i.id+'\')">Cancel</button></div>'+
  '</div></div>'
}

async function saveEditIssue(id){
  const labels=(document.getElementById('ei-labels').value||'').split(',').map(s=>s.trim()).filter(Boolean);
  const body={
    title:document.getElementById('ei-title').value,
    body:document.getElementById('ei-body').value,
    priority:document.getElementById('ei-pri').value,
    assignee:document.getElementById('ei-assign').value,
    labels:labels,
    time_estimate:parseInt(document.getElementById('ei-est').value)||0,
    time_spent:parseInt(document.getElementById('ei-spent').value)||0
  };
  await api('/api/issues/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  closeModal();loadIssues();loadFilters();
}

// ── Projects ──

function showNewProject(){
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>New Project</h2>'+
    '<label class="fl">Name</label><input type="text" id="np-name">'+
    '<label class="fl">Slug</label><input type="text" id="np-slug" placeholder="auto-generated">'+
    '<label class="fl">Description</label><textarea id="np-desc" rows="2"></textarea>'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveNewProject()">Create</button><button class="btn btn-d" onclick="closeModal()">Cancel</button></div>'+
  '</div></div>'
}

async function saveNewProject(){
  const body={name:document.getElementById('np-name').value,slug:document.getElementById('np-slug').value,description:document.getElementById('np-desc').value};
  if(!body.name){alert('Name required');return}
  const r=await api('/api/projects',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  if(r.error){alert(r.error);return}
  curProject=r.id;closeModal();init()
}

// ── Milestones ──

async function loadMilestones(){
  if(!curProject){document.getElementById('msList').innerHTML='<div class="empty">Select a project to view milestones.</div>';return}
  const d=await api('/api/projects/'+curProject+'/milestones');milestones=d.milestones||[];
  renderMilestones()
}

function renderMilestones(){
  const el=document.getElementById('msList');
  if(!milestones.length){el.innerHTML='<div class="empty">No milestones yet.</div>';return}
  el.innerHTML=milestones.map(m=>{
    const total=m.open_count+m.closed_count;
    const pct=total?Math.round(m.closed_count/total*100):0;
    const due=m.due_date?'Due: '+m.due_date:'No due date';
    return '<div class="ms-card">'+
      '<div style="display:flex;justify-content:space-between;align-items:center">'+
        '<div class="ms-title">'+esc(m.title)+' <span style="font-size:.65rem;color:'+(m.status==='open'?'var(--green)':'var(--cm)')+'">'+m.status+'</span></div>'+
        '<div style="display:flex;gap:.3rem">'+
          '<button class="btn btn-d" style="font-size:.6rem;padding:.15rem .4rem" onclick="editMs(\''+m.id+'\')">Edit</button>'+
          '<button class="btn btn-d" style="font-size:.6rem;padding:.15rem .4rem" onclick="if(confirm(\'Delete?\'))delMs(\''+m.id+'\')">Del</button>'+
        '</div>'+
      '</div>'+
      '<div class="ms-meta">'+due+' &middot; '+m.closed_count+'/'+total+' issues closed ('+pct+'%)</div>'+
      '<div class="ms-bar"><div class="ms-fill" style="width:'+pct+'%"></div></div>'+
    '</div>'
  }).join('')
}

function showNewMilestone(){
  if(!curProject){alert('Select a project first');return}
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>New Milestone</h2>'+
    '<label class="fl">Title</label><input type="text" id="nm-title">'+
    '<label class="fl">Description</label><textarea id="nm-desc" rows="2"></textarea>'+
    '<label class="fl">Due date</label><input type="date" id="nm-due">'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveNewMs()">Create</button><button class="btn btn-d" onclick="closeModal()">Cancel</button></div>'+
  '</div></div>'
}

async function saveNewMs(){
  const body={project_id:curProject,title:document.getElementById('nm-title').value,description:document.getElementById('nm-desc').value,due_date:document.getElementById('nm-due').value};
  if(!body.title){alert('Title required');return}
  await api('/api/milestones',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  closeModal();loadMilestones()
}

async function editMs(id){
  const m=milestones.find(x=>x.id===id);if(!m)return;
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>Edit Milestone</h2>'+
    '<label class="fl">Title</label><input type="text" id="em-title" value="'+esc(m.title)+'">'+
    '<label class="fl">Description</label><textarea id="em-desc" rows="2">'+esc(m.description)+'</textarea>'+
    '<label class="fl">Due date</label><input type="date" id="em-due" value="'+esc(m.due_date)+'">'+
    '<label class="fl">Status</label><select id="em-status"><option value="open"'+(m.status==='open'?' selected':'')+'>Open</option><option value="closed"'+(m.status==='closed'?' selected':'')+'>Closed</option></select>'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveEditMs(\''+id+'\')">Save</button><button class="btn btn-d" onclick="closeModal()">Cancel</button></div>'+
  '</div></div>'
}

async function saveEditMs(id){
  const body={title:document.getElementById('em-title').value,description:document.getElementById('em-desc').value,due_date:document.getElementById('em-due').value,status:document.getElementById('em-status').value};
  await api('/api/milestones/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
  closeModal();loadMilestones()
}

async function delMs(id){await api('/api/milestones/'+id,{method:'DELETE'});loadMilestones()}

// ── Activity ──

async function loadActivity(){
  const d=await api('/api/activity?project_id='+(curProject||'')+'&limit=50');
  const acts=d.activity||[];
  const el=document.getElementById('actList');
  if(!acts.length){el.innerHTML='<div class="empty">No activity yet.</div>';return}
  el.innerHTML=acts.map(a=>'<div class="act-item"><span class="act-action">'+esc(a.action)+'</span> '+esc(a.detail)+'<span class="act-time">'+timeAgo(a.created_at)+'</span></div>').join('')
}

// ── Tabs ──

function switchTab(tab){
  document.querySelectorAll('.tab').forEach(t=>t.classList.toggle('active',t.dataset.tab===tab));
  document.getElementById('pane-issues').style.display=tab==='issues'?'':'none';
  document.getElementById('pane-board').style.display=tab==='board'?'':'none';
  document.getElementById('pane-milestones').style.display=tab==='milestones'?'':'none';
  document.getElementById('pane-activity').style.display=tab==='activity'?'':'none';
  if(tab==='board')loadBoard();
  if(tab==='milestones')loadMilestones();
  if(tab==='activity')loadActivity();
}

// ── Board ──

let boardIssues=[];
async function loadBoard(){
  const p=new URLSearchParams();
  if(curProject)p.set('project_id',curProject);
  p.set('status','all');p.set('limit','200');
  const d=await api('/api/issues?'+p);
  boardIssues=d.issues||[];
  renderBoard();
}

function renderBoard(){
  const cols={open:[],in_progress:[],closed:[]};
  boardIssues.forEach(i=>{
    const s=cols[i.status]!==undefined?i.status:'open';
    cols[s].push(i);
  });
  ['open','in_progress','closed'].forEach(status=>{
    const el=document.getElementById('col-'+status);
    el.innerHTML=cols[status].map(i=>{
      const priCls='pri-'+i.priority;
      const labels=(i.labels||[]).map(l=>'<span class="label-chip">'+esc(l)+'</span>').join(' ');
      return '<div class="board-card" draggable="true" data-id="'+i.id+'" ondragstart="boardDragStart(event)" ondragend="boardDragEnd(event)" onclick="showIssue(\''+i.id+'\')">'+
        '<div class="bc-title">'+esc(i.title)+'</div>'+
        '<div class="bc-meta">'+
          '<span class="bc-num">#'+i.number+'</span>'+
          '<span class="priority-badge '+priCls+'">'+i.priority+'</span>'+
          labels+
          (i.assignee?'<span>'+esc(i.assignee)+'</span>':'')+
        '</div>'+
      '</div>'
    }).join('');
  });
  document.getElementById('bc-open').textContent=cols.open.length;
  document.getElementById('bc-wip').textContent=cols.in_progress.length;
  document.getElementById('bc-closed').textContent=cols.closed.length;
}

let dragId='';
function boardDragStart(e){
  dragId=e.target.dataset.id;
  e.target.classList.add('dragging');
  e.dataTransfer.effectAllowed='move';
}
function boardDragEnd(e){e.target.classList.remove('dragging')}
function boardDragOver(e){
  e.preventDefault();e.dataTransfer.dropEffect='move';
  e.currentTarget.classList.add('drag-over');
}
function boardDragLeave(e){e.currentTarget.classList.remove('drag-over')}
async function boardDrop(e){
  e.preventDefault();
  e.currentTarget.classList.remove('drag-over');
  const newStatus=e.currentTarget.dataset.status;
  if(!dragId||!newStatus)return;
  await api('/api/issues/'+dragId+'/status',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({status:newStatus})});
  dragId='';
  loadBoard();loadStats();
}

function closeModal(){document.getElementById('modal').innerHTML=''}

init();
fetch('/api/tier').then(r=>r.json()).then(j=>{if(j.tier==='free'){document.getElementById('upgrade-banner').style.display='block';var b=document.getElementById('tier-badge');if(b){b.className='badge badge-free';b.textContent='Free'}}else{var b=document.getElementById('tier-badge');if(b){b.className='badge badge-pro';b.textContent='Pro'}}}).catch(()=>{var b=document.getElementById('upgrade-banner');if(b)b.style.display='block'});
</script><script>
(function(){
  fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
    if(!cfg||typeof cfg!=='object')return;
    if(cfg.dashboard_title){
      document.title=cfg.dashboard_title;
      var h1=document.querySelector('h1');
      if(h1){
        var inner=h1.innerHTML;
        var firstSpan=inner.match(/<span[^>]*>[^<]*<\/span>/);
        if(firstSpan){h1.innerHTML=firstSpan[0]+' '+cfg.dashboard_title}
        else{h1.textContent=cfg.dashboard_title}
      }
    }
  }).catch(function(){});
})();
</script>
</body></html>`
