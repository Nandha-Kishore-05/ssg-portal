import "./style.css";

function TextArea(props) {
  return (
    <div style={{ marginTop: props.margin }} className="input-box">
      <label>{props.label}</label>
      <br />
      <textarea
        style={{ minHeight: props.height === undefined ? 200 : props.height }}
        placeholder={props.placeholder}
        onChange={(e) => {
          props.onChange(e.target.value);
        }}
      />
    </div>
  );
}

export default TextArea;
