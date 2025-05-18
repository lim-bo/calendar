#ifndef ALLPARTICIPANTS_H
#define ALLPARTICIPANTS_H

#include <QDialog>
#include "client.h"
#include "eventData.h"
#include "event.h"

namespace Ui {
class allparticipants;
}

class allparticipants : public QDialog
{
    Q_OBJECT

public:
    explicit allparticipants(const EventData &eventData, QWidget *parent = nullptr);
    ~allparticipants();

private:
    Ui::allparticipants *ui;
    EventData eventData;
};

#endif // ALLPARTICIPANTS_H
