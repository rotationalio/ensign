import { CardListItem } from "@/components/common/CardListItem"
import { Card } from "@rotational/beacon-core"

export function InviteCard() {
    return(
        <>
      <Card>
        <Card.Header>
          <h1 className="text-base font-bold">You've Been Invited!</h1>
        </Card.Header>
        
        {/* TODO: Add conditional to card body below appears on new user invite page */}
        <Card.Body>
            <p>You've been invited by <span className="font-bold">(inviter name)</span> to join the <span className="font-bold">(org name)</span> organization as <span className="font-bold">(role)</span> on Ensign! Create your account today.</p>
        </Card.Body>

        {/* TODO: Add conditional so that card body below appears on existing user invite page */}
        <Card.Body>
            <p>You've been invited by <span className="font-bold">(inviter name)</span> to join the <span className="font-bold">(org name)</span> organization as <span className="font-bold">(role)</span> on Ensign! Log in to accept the invitation.</p>
        </Card.Body>

      </Card>
        </>
    )
}

export default InviteCard