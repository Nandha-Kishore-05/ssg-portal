import "./style.css";

function CustomButton(props) {
  return (
    <button
      onClick={props.onClick}
      className="custom-button"
      style={{
        width: props.width ? props.width : "100%",
        marginTop: props.margin,
        backgroundColor:
          props.type === "success"
            ? "#605BFF"
            : props.type === "danger"
            ? "red"
            : "",
      }}
    >
      {props.label}
    </button>
  );
}

export default CustomButton;
