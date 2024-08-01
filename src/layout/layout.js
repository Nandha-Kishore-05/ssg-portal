import SideBar from "./sidebar"

function AppLayout(props){
    return (
        <div className="app-layout">
            <SideBar rId={props.rId} />
            <div className="app-body">
                <h1>{props.title}</h1>
                {props.body}
            </div>
        </div>
    )
}

export default AppLayout