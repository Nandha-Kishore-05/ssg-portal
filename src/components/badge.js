import "./style.css";

function Badge(props) {
  return (
    <div
      style={{
        width:"max-content",
        color:
          props.type === "success"
            ? "green"
            : props.type === "danger"
            ? "red"
            : "",
        backgroundColor:
          props.type === "success"
            ?"#00800030"
            : props.type === "danger"
            ? "#ff00002e"
            : "",
      }}
      className="app-badge"
    >
      <h5>{props.label}</h5>
    </div>
  );
}

export default Badge;
