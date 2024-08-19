import "./style.css";

function InputBox(props) {
  return (
    <div style={{ marginTop: props.margin }} className="input-box">
      <label>{props.label}</label>
      <br />
      <input
        placeholder={props.placeholder}
        type={props.type}
        accept={props.accept}
        value={props.value}
        onChange={(e) => {
          props.type === "file"
            ? props.onChange(e)
            : props.onChange(e.target.value);
        }}
      />
    </div>
  );
}

export default InputBox;
