import "./style.css";

function InputBox(props) {
  return (
    <div className="input-box" style={{ marginTop: props.margin }}>
      <label>{props.label}</label>
      <br />
      <input type={props.type} placeholder={props.placeholder} />
    </div>
  );
}

export default InputBox;
