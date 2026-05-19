# Events

Events define message payloads for asynchronous communication or pub/sub architectures. 

## Defining an Event

Use the `event` keyword followed by the event name and a block specifying its properties.

```contractor
event SignUpEvent {
    name: "Sign Up"
    payload: User<String>
}
```

## Event Properties

Events support the following properties:
- `name`: (Optional) The string identifier for the event. If not provided, the name of the event block (`SignUpEvent`) is typically used.
- `payload`: The type of data associated with the event. This can be a primitive type or a defined `model`.
