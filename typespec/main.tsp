import "@typespec/http";

using Http;
@service(#{ title: "Report-Engine Server Service" })
// TODO make this configurable
@server("http://localhost:8080", "Report-Engine API Server")
@route("/api")
namespace ServerService;

model PingResponse {
  message: string;
}

// @error
// model Error {
//   code: int32;
//   message: string;
// }

@route("/ping")
@tag("ping")
interface ping {
  @get op pong(): PingResponse;
}
